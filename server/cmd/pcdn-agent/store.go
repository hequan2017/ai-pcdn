package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store 本地 JSONL 持久化（pending 队列）。每行一个 TrafficPoint。
// 上报成功后整体清空；失败保留由重试任务兜底。并发由 mu 保护。
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

// ReadAll 读取全部 pending 点
func (s *Store) ReadAll() ([]TrafficPoint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var points []TrafficPoint
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		var p TrafficPoint
		if json.Unmarshal([]byte(line), &p) == nil {
			points = append(points, p)
		}
	}
	return points, nil
}

// Clear 清空 pending（上报成功后调用）
func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return os.Truncate(s.path, 0)
}
