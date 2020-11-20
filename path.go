package main

import (
	"errors"
	"strings"
)

// ErrBadPath is returned when a bad UNC share path is provided.
var ErrBadPath = errors.New(`UNC share path should be of the form \\<server>\<share>`)

// ServerAddrFromPath returns the server address for the given UNC path
func ServerAddrFromPath(p string) (string, error) {
	// Remove prefix
	if !strings.HasPrefix(p, `\\`) {
		return "", ErrBadPath
	}
	p = strings.TrimPrefix(p, `\\`)

	// Remove share name
	parts := strings.SplitN(p, `\`, 2)
	if len(parts) != 2 {
		return "", ErrBadPath
	}
	p = parts[0]

	// Add port number if one isn't provided
	if !strings.Contains(p, `:`) {
		p += ":445"
	}

	return p, nil
}
