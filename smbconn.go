package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/hirochachacha/go-smb2"
)

// ErrClosed is returned when an operation is attempted on a closed connection.
var ErrClosed = errors.New("smbfs: the connection to the remote file system is closed")

// SMBConn is an individual
type SMBConn struct {
	source string
	id     int

	mutex   sync.RWMutex
	session *smb2.Session
	share   *smb2.Share
}

// NewConnection establishes a new SMB connection to a remote file system.
func NewConnection(ctx context.Context, conf Config, id int) (*SMBConn, error) {
	host, err := ServerAddrFromPath(conf.Source)
	if err != nil {
		return nil, fmt.Errorf("smbfs: unable to determine source host: %v", err)
	}

	// Prepare a connetion context to limit the connection time
	connCtx, cancel := context.WithTimeout(ctx, conf.ConnTimeout)
	defer cancel()

	// Prepare a TCP dialer
	var tcpd net.Dialer

	// Establish a TCP connection via the dialer
	conn, err := tcpd.DialContext(connCtx, "tcp", host)
	if err != nil {
		return nil, err
	}

	// Prepare an SMB2 dialer
	smbd := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     conf.User,
			Password: conf.Password,
			Domain:   conf.Domain,
		},
	}

	// Establish an SMB2 session via the dialer
	session, err := smbd.DialContext(connCtx, conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	session = session.WithContext(ctx)

	// Mount the desired path
	share, err := session.Mount(conf.Source)
	if err != nil {
		session.Logoff()
		return nil, err
	}
	share = share.WithContext(ctx)

	c := &SMBConn{
		id:      id,
		source:  conf.Source,
		session: session,
		share:   share,
	}

	c.log("Connected")

	return c, nil
}

// Close disconnects from the SMB server.
func (c *SMBConn) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Make sure the connection hasn't been closed already
	if c.session == nil {
		return nil
	}

	// Execute the logoff process
	err := c.session.Logoff()

	// Change state to disconnected
	c.session = nil
	c.share = nil

	c.log("Closed")

	return err
}

// Open returns a file from the remote file sytem.
func (c *SMBConn) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, `/`) // Remove leading slash

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.share == nil {
		return nil, ErrClosed
	}

	f, err := c.share.Open(name)
	if err != nil {
		switch err := err.(type) {
		case *os.PathError:
			c.log("%s: %v", name, err.Unwrap())
		default:
			c.log("%v", err)
		}
	} else {
		c.log("%s", name)
	}
	return f, err
}

// OK returns true if the connection is healthy.
func (c *SMBConn) OK() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.share == nil {
		return false
	}

	// If we can perform a stat operation on the root of the share it's a
	// pretty good indication that the connection is healthy
	_, err := c.share.Stat("")
	if err != nil {
		switch err := err.(type) {
		case *os.PathError:
			c.log("HEALTH CHECK: %v", err.Unwrap())
		default:
			c.log("HEALTH CHECK: %v", err)
		}
		return false
	}

	c.log("HEALTH OK")

	return true
}

func (c *SMBConn) log(format string, v ...interface{}) {
	fmt.Printf("[%s][%d]: %s\n", c.source, c.id, fmt.Sprintf(format, v...))
}
