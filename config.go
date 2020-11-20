package main

import (
	"fmt"
	"time"
)

// Config holds a set of SMB HTTP proxy configuration values.
type Config struct {
	URLPrefix   string        `kong:"env='URL_PREFIX',required"`
	Source      string        `kong:"env='SOURCE',required"`
	User        string        `kong:"env='USER',required"`
	Password    string        `kong:"env='PASSWORD',required"`
	Domain      string        `kong:"env='DOMAIN',required"`
	ConnTimeout time.Duration `kong:"env='CONN_TIMEOUT',required,default='5s'"`
}

// Summary returns a multiline string representation of the configuration.
func (c Config) Summary() string {
	host, _ := ServerAddrFromPath(c.Source)
	return fmt.Sprintf("URL Prefix: %s\nSource: %s\nHost: %s\nUser: %s\nDomain: %s", c.URLPrefix, c.Source, host, c.User, c.Domain)
}
