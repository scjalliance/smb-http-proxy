// +build !linux

package main

import "errors"

// DefaultMountFlags is a default set of readonly mount flags.
const DefaultMountFlags = 0

// DefaultUnmountFlags is a default set of forced unmount flags.
const DefaultUnmountFlags = 0

func mount(source string, target string, fstype string, flags uintptr, data string) error {
	return errors.New("mount is unsupported on non-linux platforms")
}

func unmount(target string, flags int) error {
	return errors.New("unmount is unsupported on non-linux platforms")
}
