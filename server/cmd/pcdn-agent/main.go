// pcdn-agent PCDN 采集 agent：部署到每台节点服务器，采集网卡流量分钟峰值并上报管理后台。
//
// 用法：
//   pcdn-agent -server https://api.example.com -sn PCDN-xxxx -token xxxx
//
// 数据流：每秒采样网卡速率 → 每分钟取峰值写入本地 pending 文件 → 上报 → 失败保留由独立重试任务兜底。
// 心跳 30s，首次启动激活。上报幂等由后台 (node_id, window_start, iface) 唯一索引保证。
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	var (
		server   = flag.String("server", "http://127.0.0.1:8888", "管理后台地址(含协议与前缀)")
		sn       = flag.String("sn", "", "节点序列号")
		token    = flag.String("token", "", "节点 token")
		ifaces   = flag.String("ifaces", "", "采集网卡(逗号分隔,空=所有非lo)")
		interval = flag.Int("interval", 60, "采集与上报间隔(秒)")
		storePath = flag.String("store", "/var/lib/pcdn-agent/pending.jsonl", "本地持久化路径")
	)
	flag.Parse()
	if *sn == "" || *token == "" {
		log.Fatal("必须提供 -sn 和 -token")
	}

	dur := time.Duration(*interval) * time.Second
	store := NewStore(*storePath)
	reporter := NewReporter(*server, *sn, *token)
	collector := NewCollector(parseIfaces(*ifaces))

	// 首次激活（失败不阻塞，心跳会持续刷新状态）
	if err := reporter.Activate(getHostInfo()); err != nil {
		log.Printf("激活失败(将在心跳中重试): %v", err)
	}

	go collectorLoop(collector, store, dur)
	go reporterLoop(reporter, store, dur)
	go retryLoop(reporter, store, dur*2) // 独立重试任务，与即时上报解耦
	go heartbeatLoop(reporter, 30*time.Second)

	log.Printf("pcdn-agent 已启动 sn=%s server=%s ifaces=%v", *sn, *server, *ifaces)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Printf("pcdn-agent 退出")
}

func parseIfaces(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

// collectorLoop 每秒采样速率，每 interval 取峰值落盘
func collectorLoop(c *Collector, s *Store, interval time.Duration) {
	sec := time.NewTicker(time.Second)
	min := time.NewTicker(interval)
	defer sec.Stop()
	defer min.Stop()
	for {
		select {
		case <-sec.C:
			c.Sample()
		case <-min.C:
			for _, p := range c.FlushMinute() {
				if err := s.Append(p); err != nil {
					log.Printf("本地持久化失败: %v", err)
				}
			}
		}
	}
}

// reporterLoop 每 interval 上报 pending 点，失败保留
func reporterLoop(r *Reporter, s *Store, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for range t.C {
		reportPending(r, s)
	}
}

// retryLoop 独立兜底重试，扫 pending 重报
func retryLoop(r *Reporter, s *Store, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for range t.C {
		reportPending(r, s)
	}
}

// heartbeatLoop 每 interval 发心跳
func heartbeatLoop(r *Reporter, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for range t.C {
		if err := r.Heartbeat(getHostInfo()); err != nil {
			log.Printf("心跳失败: %v", err)
		}
	}
}

// reportPending 读全部 pending → 上报 → 成功清空（reporter 内部互斥锁串行化，避免并发上报）
func reportPending(r *Reporter, s *Store) {
	if err := r.ReportPending(s); err != nil {
		log.Printf("上报失败: %v", err)
	}
}
