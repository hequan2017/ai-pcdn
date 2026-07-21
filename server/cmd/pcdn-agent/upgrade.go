package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	if info.Checksum != "" {
		sum, err := sha256File(tmp)
		if err != nil {
			os.Remove(tmp)
			return err
		}
		if sum != info.Checksum {
			os.Remove(tmp)
			return fmt.Errorf("checksum mismatch")
		}
	}
	if err := os.Chmod(tmp, 0o755); err != nil {
		os.Remove(tmp)
		return err
	}
	old := os.Args[0] + ".old"
	os.Remove(old)
	if err := os.Rename(os.Args[0], old); err != nil {
		return err
	}
	if err := os.Rename(tmp, os.Args[0]); err != nil {
		os.Rename(old, os.Args[0])
		return err
	}
	log.Printf("升级完成，重启")
	return doRestart()
}

func downloadFile(url, path string) error {
	resp, err := http.Get(url)
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
	_, err = io.Copy(out, resp.Body)
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
