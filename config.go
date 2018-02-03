package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gentlemanautomaton/bindflag"
)

// Config holds a set of SMB HTTP proxy configuration values.
type Config struct {
	URLPrefix string
	Source    string
	Target    string
	Username  string
	Password  string
	Domain    string
	Options   string
}

// DefaultConfig holds the default configuration values.
var DefaultConfig = Config{
	Target:  "/mnt/smb",
	Options: "ro,uid=0,gid=0,forceuid,forcegid,vers=2.1,sec=ntlm",
}

// ParseEnv will parse environment variables and apply them to the
// configuration.
func (c *Config) ParseEnv() {
	var (
		prefix, hasPrefix     = os.LookupEnv("URLPREFIX")
		source, hasSource     = os.LookupEnv("UNCPATH")
		target, hasTarget     = os.LookupEnv("TARGET")
		username, hasUsername = os.LookupEnv("USERNAME")
		password, hasPassword = os.LookupEnv("PASSWORD")
		domain, hasDomain     = os.LookupEnv("DOMAIN")
		options, hasOptions   = os.LookupEnv("OPTIONS")
	)

	if hasPrefix {
		c.URLPrefix = prefix
	}
	if hasSource {
		c.Source = source
	}
	if hasTarget {
		c.Target = target
	}
	if hasUsername {
		c.Username = username
	}
	if hasPassword {
		c.Password = password
	}
	if hasDomain {
		c.Domain = domain
	}
	if hasOptions {
		c.Options = options
	}
}

// MountFlags returns the flags for syscall.Mount
func (c *Config) MountFlags() uintptr {
	return DefaultMountFlags
}

// MountOptions returns the options string for syscall.Mount
func (c *Config) MountOptions() string {
	var options []string
	if c.Username != "" {
		options = append(options, c.option("username", c.Username))
	}
	if c.Password != "" {
		options = append(options, c.option("password", c.Password))
	}
	if c.Domain != "" {
		options = append(options, c.option("domain", c.Domain))
	}
	if c.Options != "" {
		options = append(options, c.Options)
	}
	return strings.Join(options, ",")
}

// ParseArgs parses the given argument list and applies them to the
// configuration.
func (c *Config) ParseArgs(args []string, errorHandling flag.ErrorHandling) error {
	fs := flag.NewFlagSet("", errorHandling)
	c.Bind(fs)
	return fs.Parse(args)
}

// Bind will bind the given flag set to the configuration.
func (c *Config) Bind(fs *flag.FlagSet) {
	fs.Var(bindflag.String(&c.URLPrefix), "urlprefix", "URL Prefix")
	fs.Var(bindflag.String(&c.Source), "uncpath", "source UNC path (required)")
	fs.Var(bindflag.String(&c.Target), "target", "target mount path on local filesystem")
	fs.Var(bindflag.String(&c.Username), "username", "smb username (required)")
	fs.Var(bindflag.String(&c.Password), "password", "smb password (required)")
	fs.Var(bindflag.String(&c.Domain), "domain", "smb domain (optional)")
	fs.Var(bindflag.String(&c.Options), "options", "smb mount options")
}

// Validate checks the configuration and exits the program if it's invalid.
func (c *Config) Validate() {
	c.checkRequired("UNCPATH", c.Source)
	c.checkRequired("USERNAME", c.Username)
	c.checkRequired("PASSWORD", c.Password)
}

func (c *Config) checkRequired(variable, value string) {
	if value == "" {
		fmt.Printf("%s environment variable or argument is missing.\n", variable)
		os.Exit(1)
	}
}

func (c *Config) option(name, value string) string {
	if value != "" {
		return name + "=" + value
	}
	return ""
}
