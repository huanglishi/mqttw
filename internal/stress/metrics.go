// 压测指标：全内存原子计数，压测不落库（避免 SQLite 成为瓶颈）
package stress

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// echoSampleCap 回显延迟样本缓冲上限（保留最近 N 条用于分位计算）
const echoSampleCap = 10000

// Metrics 原子指标集合（所有字段并发安全）
type Metrics struct {
	ConnectTotal   atomic.Int64
	ConnectSuccess atomic.Int64
	ConnectFail    atomic.Int64
	ConnectTimeSum atomic.Int64 // ns
	ConnectTimeMax atomic.Int64 // ns

	SubscribeTotal   atomic.Int64
	SubscribeSuccess atomic.Int64
	SubscribeFail    atomic.Int64

	PublishTotal   atomic.Int64
	PublishSuccess atomic.Int64
	PublishFail    atomic.Int64

	ReceivedTotal  atomic.Int64
	EchoCount      atomic.Int64
	EchoLatencySum atomic.Int64 // ns
	EchoLatencyMax atomic.Int64 // ns

	// 失败原因分类计数
	mu          sync.Mutex
	failReasons map[string]int64

	// 回显延迟样本缓冲（最近 10000 条，用于 P50/P95/P99）
	echoMu      sync.Mutex
	echoSamples []int64
}

func newMetrics() *Metrics {
	return &Metrics{
		failReasons:  make(map[string]int64),
		echoSamples:  make([]int64, 0, echoSampleCap),
	}
}

// AddFailReason 失败原因 +1（如 "connect timeout" / "connection refused" / "not authorized"）
func (m *Metrics) AddFailReason(reason string) {
	if reason == "" {
		reason = "unknown"
	}
	m.mu.Lock()
	m.failReasons[reason]++
	m.mu.Unlock()
}

// AddConnectTime 记录单次连接耗时
func (m *Metrics) AddConnectTime(d time.Duration) {
	ns := d.Nanoseconds()
	m.ConnectTimeSum.Add(ns)
	for {
		old := m.ConnectTimeMax.Load()
		if ns <= old || m.ConnectTimeMax.CompareAndSwap(old, ns) {
			break
		}
	}
}

// AddEchoLatency 记录一次回显延迟，并追加分位样本
func (m *Metrics) AddEchoLatency(d time.Duration) {
	ns := d.Nanoseconds()
	m.EchoCount.Add(1)
	m.EchoLatencySum.Add(ns)
	for {
		old := m.EchoLatencyMax.Load()
		if ns <= old || m.EchoLatencyMax.CompareAndSwap(old, ns) {
			break
		}
	}
	// 追加样本（保留最近 10000 条）
	m.echoMu.Lock()
	m.echoSamples = append(m.echoSamples, ns)
	if len(m.echoSamples) > echoSampleCap {
		m.echoSamples = append([]int64(nil), m.echoSamples[len(m.echoSamples)-echoSampleCap:]...)
	}
	m.echoMu.Unlock()
}

// Sample 快照（startAt=任务开始时间，now=当前时间；返回每秒速率）
type sample struct {
	ts         int64 // unix 秒，用于前端时间轴
	connected  int64 // 当前已连接数（Runner 维护）
	connectAvg int64 // 平均连接耗时 ms
	connectMax int64
	pubRate    int64 // msg/s
	recvRate   int64
	echoAvg    int64 // 回显平均延迟 ms
	echoMax    int64
}

// rawSnapshot 基础数值快照（不含速率）
type rawSnapshot struct {
	Status          string           `json:"status"`
	Connected       int64            `json:"connected"`
	ConnectTotal    int64            `json:"connect_total"`
	ConnectSuccess  int64            `json:"connect_success"`
	ConnectFail     int64            `json:"connect_fail"`
	ConnectAvgMs    int64            `json:"connect_avg_ms"`
	ConnectMaxMs    int64            `json:"connect_max_ms"`
	SubscribeTotal  int64            `json:"subscribe_total"`
	SubscribeSucc   int64            `json:"subscribe_success"`
	SubscribeFail   int64            `json:"subscribe_fail"`
	PublishTotal    int64            `json:"publish_total"`
	PublishSuccess  int64            `json:"publish_success"`
	PublishFail     int64            `json:"publish_fail"`
	ReceivedTotal   int64            `json:"received_total"`
	EchoCount       int64            `json:"echo_count"`
	EchoAvgMs       int64            `json:"echo_avg_ms"`
	EchoMaxMs       int64            `json:"echo_max_ms"`
	EchoP50Ms       int64            `json:"echo_p50_ms"`
	EchoP95Ms       int64            `json:"echo_p95_ms"`
	EchoP99Ms       int64            `json:"echo_p99_ms"`
	FailReasons     map[string]int64 `json:"fail_reasons"`
	StartAt         int64            `json:"start_at"`     // 任务开始时间 unix 毫秒
	DurationSec     int64            `json:"duration_sec"` // 计划时长（秒），0=不限
	Series          []map[string]any `json:"series"` // 时间序列（前端折线图）
}

// buildRaw 组装当前数值快照
func (m *Metrics) buildRaw(status string, connected int64) rawSnapshot {
	snap := rawSnapshot{
		Status:         status,
		Connected:      connected,
		ConnectTotal:   m.ConnectTotal.Load(),
		ConnectSuccess: m.ConnectSuccess.Load(),
		ConnectFail:    m.ConnectFail.Load(),
		SubscribeTotal: m.SubscribeTotal.Load(),
		SubscribeSucc:  m.SubscribeSuccess.Load(),
		SubscribeFail:  m.SubscribeFail.Load(),
		PublishTotal:   m.PublishTotal.Load(),
		PublishSuccess: m.PublishSuccess.Load(),
		PublishFail:    m.PublishFail.Load(),
		ReceivedTotal:  m.ReceivedTotal.Load(),
		EchoCount:      m.EchoCount.Load(),
		FailReasons:    map[string]int64{},
	}
	// 连接耗时均值（ms）
	if n := snap.ConnectSuccess; n > 0 {
		snap.ConnectAvgMs = m.ConnectTimeSum.Load() / n / int64(time.Millisecond)
	}
	snap.ConnectMaxMs = m.ConnectTimeMax.Load() / int64(time.Millisecond)
	// 回显延迟均值（ms）
	if n := snap.EchoCount; n > 0 {
		snap.EchoAvgMs = m.EchoLatencySum.Load() / n / int64(time.Millisecond)
	}
	snap.EchoMaxMs = m.EchoLatencyMax.Load() / int64(time.Millisecond)
	// 回显延迟分位 P50/P95/P99（样本拷贝排序）
	m.echoMu.Lock()
	samples := make([]int64, len(m.echoSamples))
	copy(samples, m.echoSamples)
	m.echoMu.Unlock()
	if len(samples) > 0 {
		sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
		snap.EchoP50Ms = percentile(samples, 0.50) / int64(time.Millisecond)
		snap.EchoP95Ms = percentile(samples, 0.95) / int64(time.Millisecond)
		snap.EchoP99Ms = percentile(samples, 0.99) / int64(time.Millisecond)
	}

	m.mu.Lock()
	for k, v := range m.failReasons {
		snap.FailReasons[k] = v
	}
	m.mu.Unlock()
	return snap
}

// percentile 从已排序样本计算 p 分位（0-1）
func percentile(sorted []int64, p float64) int64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p)
	if idx < 0 {
		idx = 0
	}
	return sorted[idx]
}

// sortFailReasons 失败原因降序排列（用于前端展示 Top N）
func sortFailReasons(m map[string]int64) [][2]any {
	out := make([][2]any, 0, len(m))
	for k, v := range m {
		out = append(out, [2]any{k, v})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i][1].(int64) > out[j][1].(int64)
	})
	return out
}
