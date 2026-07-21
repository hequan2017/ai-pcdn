package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store 本地 JSONL 持久化（pending 队列）。每行一个 TrafficPoint。
// 并发由 mu 保护；上报流程用 ReadAll(返回已读边界 offset) + ClearTo(offset)，
// 只截断已读部分，保留读取期间新 Append 的点，避免数据丢失。
type Store struct {
	mu   sync.Mutex
	path string
}

func NewStore(path string) *Store {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	return &Store{path: path}
}

// Append 追加一个流量点
func (s *Store) Append(p TrafficPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = f.Write(b)
	return err
}

// ReadAll 读取全部 pending 点，并返回读取时的文件大小（已读边界 offset）。
// 调用方上报成功后应调 ClearTo(offset) 只清理已读部分。
func (s *Store) ReadAll() (points []TrafficPoint, offset int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, statErr := os.Stat(s.path)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			return nil, 0, nil
		}
		return nil, 0, statErr
	}
	offset = info.Size()
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, offset, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		var p TrafficPoint
		if json.Unmarshal([]byte(line), &p) == nil {
			points = append(points, p)
		} else {
			log.Printf("跳过损坏的 pending 行: %q", line)
		}
	}
	return points, offset, nil
}

// ClearTo 截断已读部分：若期间有新数据追加（当前 size > offset），保留 [offset, 当前) 并移到文件头；否则清空。
func (s *Store) ClearTo(offset int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, err := os.Stat(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	curSize := info.Size()
	if offset <= 0 || offset >= curSize {
		return os.Truncate(s.path, 0)
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		return os.Truncate(s.path, 0)
	}
	return os.WriteFile(s.path, data[offset:], 0o644)
}
