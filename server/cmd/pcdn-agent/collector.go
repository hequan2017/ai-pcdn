package main

import (
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Collector 网卡流量采集器：按实际时间差计算 bps，每分钟取峰值
type Collector struct {
	ifaces   map[string]bool        // 空=所有非lo
	last     map[string][2]uint64   // 上次字节计数
	lastTime map[string]time.Time   // 上次采样时间
	curMax   map[string][2]int64    // 当前分钟峰值 [rxBps, txBps]
}

func NewCollector(ifaces []string) *Collector {
	c := &Collector{
		ifaces:   map[string]bool{},
		last:     map[string][2]uint64{},
		lastTime: map[string]time.Time{},
		curMax:   map[string][2]int64{},
	}
	for _, f := range ifaces {
		c.ifaces[strings.TrimSpace(f)] = true
	}
	return c
}

// Sample 采样：按实际时间差计算 bps（避免 ticker 漂移导致 bps 高估），更新当前分钟峰值
func (c *Collector) Sample() {
	cur, err := readProcNetDev()
	if err != nil {
		return
	}
	now := time.Now()
	for iface, bytes := range cur {
		if iface == "lo" {
			continue
		}
		if len(c.ifaces) > 0 && !c.ifaces[iface] {
			continue
		}
		last, ok := c.last[iface]
		if !ok {
			c.last[iface] = bytes
			c.lastTime[iface] = now
			continue
		}
		dt := now.Sub(c.lastTime[iface]).Seconds()
		c.last[iface] = bytes
		c.lastTime[iface] = now
		if dt <= 0 {
			continue
		}
		rxBps := int64(float64(bytes[0]-last[0]) * 8 / dt)
		txBps := int64(float64(bytes[1]-last[1]) * 8 / dt)
		mx := c.curMax[iface]
		if rxBps > mx[0] {
			mx[0] = rxBps
		}
		if txBps > mx[1] {
			mx[1] = txBps
		}
		c.curMax[iface] = mx
	}
}

// FlushMinute 取出当前分钟峰值并重置，窗口对齐到分钟
func (c *Collector) FlushMinute() []TrafficPoint {
	windowStart := time.Now().Truncate(time.Minute)
	points := make([]TrafficPoint, 0, len(c.curMax))
	for iface, mx := range c.curMax {
		points = append(points, TrafficPoint{
			IfaceName:   iface,
			WindowStart: windowStart,
			RxMaxBps:    mx[0],
			TxMaxBps:    mx[1],
		})
	}
	c.curMax = map[string][2]int64{}
	return points
}

// readProcNetDev 解析 /proc/net/dev，返回 map[iface]{rxBytes, txBytes}
func readProcNetDev() (map[string][2]uint64, error) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	result := map[string][2]uint64{}
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		iface := strings.TrimSpace(parts[0])
		fields := strings.Fields(strings.TrimSpace(parts[1]))
		if len(fields) < 16 {
			continue
		}
		rx, _ := strconv.ParseUint(fields[0], 10, 64)
		tx, _ := strconv.ParseUint(fields[8], 10, 64)
		result[iface] = [2]uint64{rx, tx}
	}
	return result, nil
}

// getHostInfo 采集本机硬件信息（激活/心跳上报）
func getHostInfo() HostInfo {
	hostname, _ := os.Hostname()
	return HostInfo{
		Hostname: hostname,
		OS:       runtime.GOOS,
		InnerIP:  getOutboundIP(),
		Ifaces:   listIfaceInfo(),
	}
}

func listIfaceInfo() []IfaceInfo {
	var result []IfaceInfo
	ifaces, err := net.Interfaces()
	if err != nil {
		return result
	}
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		result = append(result, IfaceInfo{
			IfaceName: ifc.Name,
			Mac:       ifc.HardwareAddr.String(),
			Enabled:   ifc.Flags&net.FlagUp != 0,
		})
	}
	return result
}

// getOutboundIP 通过 UDP 拨号获取本机出口 IP（不实际发包）
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return addr.IP.String()
	}
	return ""
}
