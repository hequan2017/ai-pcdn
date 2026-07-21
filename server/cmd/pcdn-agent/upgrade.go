package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

// VersionInfo 后端返回的最新版本信息
type VersionInfo struct {
	Latest      bool   `json:"latest"`
	Version     string `json:"version"`
	DownloadURL string `json:"downloadUrl"`
	Checksum    string `json:"checksum"`
	Force       bool   `json:"force"`
}

// 下载客户端：限制超时，防止慢源阻塞 upgradeLoop
var dlClient = &http.Client{Timeout: 5 * time.Minute}

// maxDownloadBytes 下载大小上限，防止恶意大文件撑爆磁盘
const maxDownloadBytes = 500 * (1 << 20)

// upgradeLoop 每小时检查自升级
func upgradeLoop(r *Reporter, currentVersion string) {
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for range t.C {
		if err := CheckAndUpgrade(r, currentVersion); err != nil {
			log.Printf("升级检查失败: %v", err)
		}
	}
}

// CheckAndUpgrade 发现新版本则下载、校验、替换、重启
func CheckAndUpgrade(r *Reporter, currentVersion string) error {
	// Windows 不支持运行中 exe 替换与 syscall.Exec 重启，跳过自升级
	if runtime.GOOS == "windows" {
		return nil
	}
	info, err := r.GetVersion()
	if err != nil || !info.Latest {
		return nil
	}
	if info.Version == currentVersion && !info.Force {
		return nil
	}
	log.Printf("发现新版本 %s（当前 %s），开始升级", info.Version, currentVersion)

	tmp := os.Args[0] + ".new"
	if err := downloadFile(info.DownloadURL, tmp); err != nil {
		return err
	}
	// checksum 必填：空则拒绝升级（防止下载源被替换/投毒）
	if info.Checksum == "" {
		os.Remove(tmp)
		return fmt.Errorf("服务端未提供 checksum，拒绝升级")
	}
	sum, err := sha256File(tmp)
	if err != nil {
		os.Remove(tmp)
		return err
	}
	if sum != info.Checksum {
		os.Remove(tmp)
		return fmt.Errorf("checksum mismatch")
	}
	if err := os.Chmod(tmp, 0o755); err != nil {
		os.Remove(tmp)
		return err
	}
	old := os.Args[0] + ".old"
	os.Remove(old)
	if err := os.Rename(os.Args[0], old); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, os.Args[0]); err != nil {
		os.Rename(old, os.Args[0])
		os.Remove(tmp)
		return err
	}
	log.Printf("升级完成，重启")
	return doRestart()
}

func downloadFile(url, path string) error {
	resp, err := dlClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("download status %d", resp.StatusCode)
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, io.LimitReader(resp.Body, maxDownloadBytes))
	return err
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
