//go:build !windows

package main

import (
	"os"
	"syscall"
)

// doRestart Unix 下用 syscall.Exec 原地重启进程
func doRestart() error {
	return syscall.Exec(os.Args[0], os.Args, os.Environ())
}
