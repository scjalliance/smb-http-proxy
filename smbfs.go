package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
)

// ErrAlreadyConnected is returned when a file system is already in a connected state.
var ErrAlreadyConnected = errors.New("smbfs: a connection to the remote file system has already been established")

// SMBFS is capable of opening files on a remote file system over the SMB protocol.
type SMBFS struct {
	conf Config

	mutex  sync.RWMutex
	connID int
	conn   *SMBConn
	ctx    context.Context
}

// NewFS returns a new SMB filesystem
func NewFS(c Config) *SMBFS {
	return &SMBFS{
		conf: c,
	}
}

// Connect establishes a connection to the file server.
func (fs *SMBFS) Connect(ctx context.Context) error {
	// Hold a write lock for the duration of the call
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if fs.conn != nil {
		return ErrAlreadyConnected
	}

	// Establish an SMB connection to the remote system
	nextID := fs.connID + 1
	conn, err := NewConnection(ctx, fs.conf, nextID)
	if err != nil {
		return err
	}

	// Record the connection for use
	fs.conn = conn
	fs.connID = nextID

	// Record the context in case we need to reconnect later
	fs.ctx = ctx

	return nil
}

// Close disconnects from the SMB server.
func (fs *SMBFS) Close() error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if fs.conn == nil {
		return nil // Already closed
	}

	err := fs.conn.Close()
	fs.conn = nil
	fs.ctx = nil

	return err
}

// Open returns a file from the remote file sytem.
func (fs *SMBFS) Open(name string) (http.File, error) {
	// Attempt to open the file with the existing connection
	f, reconnect, err := fs.open(name, false)
	if err == nil || !reconnect {
		return f, err
	}

	// Attempt to reconnect
	attempted, connErr := fs.Reconnect()

	// If the reconnection failed don't bother trying again
	if !attempted || connErr != nil {
		return f, err
	}

	// Try again
	f, _, err = fs.open(name, true)

	return f, err
}

// Reconnect closes and reopens the connection to the file server.
func (fs *SMBFS) Reconnect() (attempted bool, err error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Make sure the file system hasn't been closed
	if fs.conn == nil {
		return false, ErrClosed
	}

	fs.log("Reconnecting")

	// Establish a new SMB connection to the remote system
	nextID := fs.connID + 1
	conn, err := NewConnection(fs.ctx, fs.conf, nextID)
	if err != nil {
		fs.log("Reconnection failure: %v", err)
		return true, err
	}

	// Close the old connection
	fs.conn.Close()

	// Record the new connection for use
	fs.conn = conn
	fs.connID = nextID

	fs.log("New connection established")

	return true, nil
}

func (fs *SMBFS) open(name string, retry bool) (f http.File, reconnect bool, err error) {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// Make sure the file system hasn't been closed
	if fs.conn == nil {
		return nil, false, ErrClosed
	}

	// Open the file
	f, err = fs.conn.Open(name)

	// If we succeeded, we got a permanent error, or this is a retry,
	// don't reconnect
	if err == nil || retry || isPermanentError(err) {
		return
	}

	// A reconnect is warranted if the file system connection is unhealthy
	reconnect = !fs.conn.OK()

	return
}

func (fs *SMBFS) log(format string, v ...interface{}) {
	fmt.Printf("[%s]: %s\n", fs.conf.Source, fmt.Sprintf(format, v...))
}

func isPermanentError(err error) bool {
	if os.IsNotExist(err) {
		return true
	}
	if err == context.DeadlineExceeded {
		return true
	}
	if err == context.Canceled {
		return true
	}
	return false
}
