//go:build windows

package main

import "log"

// doRestart Windows 不支持 syscall.Exec 原地重启，提示手动重启（生产环境部署在 Linux）
func doRestart() error {
	log.Printf("Windows 环境不支持自动重启，请手动重启 agent")
	return nil
}
