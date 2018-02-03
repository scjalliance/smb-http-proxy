// +build linux

package main

import "syscall"

// DefaultMountFlags is a default set of readonly mount flags.
const DefaultMountFlags = syscall.MS_RDONLY | syscall.MS_NOATIME | syscall.MS_NODIRATIME | syscall.MS_NOEXEC

// DefaultUnmountFlags is a default set of forced unmount flags.
const DefaultUnmountFlags = syscall.MNT_FORCE

func mount(source string, target string, fstype string, flags uintptr, data string) error {
	return syscall.Mount(source, target, fstype, flags, data)
}

func unmount(target string, flags int) (err error) {
	return syscall.Unmount(target, flags)
}
