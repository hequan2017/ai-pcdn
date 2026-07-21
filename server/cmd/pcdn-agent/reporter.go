package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// TrafficPoint 流量分钟峰值点（与后台 request.TrafficPoint 契约一致）
type TrafficPoint struct {
	IfaceName   string    `json:"ifaceName"`
	WindowStart time.Time `json:"windowStart"`
	RxMaxBps    int64     `json:"rxMaxBps"`
	TxMaxBps    int64     `json:"txMaxBps"`
}

type TrafficReport struct {
	Points []TrafficPoint `json:"points"`
}

type IfaceInfo struct {
	IfaceName string `json:"ifaceName"`
	Mac       string `json:"mac"`
	Enabled   bool   `json:"enabled"`
}

type HostInfo struct {
	Hostname string      `json:"hostname"`
	OS       string      `json:"os"`
	InnerIP  string      `json:"innerIp"`
	Ifaces   []IfaceInfo `json:"ifaces"`
}

type ActivateReq struct {
	Hostname string      `json:"hostname"`
	OS       string      `json:"os"`
	InnerIP  string      `json:"innerIp"`
	Ifaces   []IfaceInfo `json:"ifaces"`
}

type HeartbeatReq struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

// apiResult 后台统一响应
type apiResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Reporter 与后台通信
type Reporter struct {
	server string
	sn     string
	token  string
	client *http.Client
	mu     sync.Mutex // 串行化上报，避免 reporterLoop 与 retryLoop 并发
}

func NewReporter(server, sn, token string) *Reporter {
	return &Reporter{
		server: server,
		sn:     sn,
		token:  token,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *Reporter) post(path string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, r.server+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Node-Sn", r.sn)
	req.Header.Set("X-Node-Token", r.token)
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("http status %d", resp.StatusCode)
	}
	var res apiResult
	_ = json.NewDecoder(resp.Body).Decode(&res)
	if res.Code != 0 {
		return fmt.Errorf("server error: %s", res.Msg)
	}
	return nil
}

// Activate 首次激活
func (r *Reporter) Activate(h HostInfo) error {
	return r.post("/pcdn/agent/activate", ActivateReq{
		Hostname: h.Hostname, OS: h.OS, InnerIP: h.InnerIP, Ifaces: h.Ifaces,
	})
}

// Report 上报流量点
func (r *Reporter) Report(points []TrafficPoint) error {
	if len(points) == 0 {
		return nil
	}
	return r.post("/pcdn/agent/report", TrafficReport{Points: points})
}

// Heartbeat 心跳
func (r *Reporter) Heartbeat(h HostInfo) error {
	return r.post("/pcdn/agent/heartbeat", HeartbeatReq{Hostname: h.Hostname, OS: h.OS})
}

// GetVersion 查询最新版本（agent 自升级）
func (r *Reporter) GetVersion() (*VersionInfo, error) {
	req, err := http.NewRequest(http.MethodGet, r.server+"/pcdn/agent/version", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Node-Sn", r.sn)
	req.Header.Set("X-Node-Token", r.token)
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res struct {
		Code int         `json:"code"`
		Data VersionInfo `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	if res.Code != 0 {
		return nil, fmt.Errorf("version query failed")
	}
	return &res.Data, nil
}

// ReportPending 读 pending 全量上报，成功后清空（加锁串行化）
func (r *Reporter) ReportPending(s *Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	points, err := s.ReadAll()
	if err != nil {
		return err
	}
	if len(points) == 0 {
		return nil
	}
	if err := r.Report(points); err != nil {
		return err
	}
	return s.Clear()
}
