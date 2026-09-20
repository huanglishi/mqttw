// 压测编排：连接 → 订阅 → 发布 → 指标采样推送 → 收尾汇总
package stress

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// 运行状态
const (
	StatusIdle     = 0
	StatusRunning  = 1
	StatusStopping = 2
)

// Runner 单次压测任务编排器（同一时间只允许一个任务）
type Runner struct {
	cfg     StressConfig
	metrics *Metrics
	pool    *ClientPool
	pub     *Publisher

	status  atomic.Int32
	startAt time.Time
	stopAt  time.Time
	connOK  int64 // 本次任务成功连接数

	// 事件推送（由 service 层注入 application.Get().Event.Emit）
	emitFn func(event string, payload string)
	// 速率采样历史
	mu     sync.Mutex
	series []map[string]any
	// 上一秒累计值（算速率）
	lastPub, lastRecv int64
	lastSampleTime    time.Time

	ctx    context.Context
	cancel context.CancelFunc
}

// NewRunner 创建压测任务
func NewRunner(cfg StressConfig) *Runner {
	return &Runner{
		cfg:     cfg,
		metrics: newMetrics(),
	}
}

// SetEventEmitter 注入 Wails 事件推送
func (r *Runner) SetEventEmitter(fn func(event, payload string)) {
	r.emitFn = fn
}

// Status 当前状态
func (r *Runner) Status() int { return int(r.status.Load()) }

// Start 启动任务（非阻塞）
func (r *Runner) Start() error {
	if !r.status.CompareAndSwap(StatusIdle, StatusRunning) {
		return fmt.Errorf("已有压测任务在运行")
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.metrics = newMetrics()
	r.pool = NewClientPool(r.cfg, r.metrics)
	if r.cfg.PublishEnabled {
		r.pub = NewPublisher(r.cfg, r.metrics, r.pool)
	}
	r.startAt = time.Now()
	r.lastSampleTime = r.startAt
	r.lastPub, r.lastRecv = 0, 0
	r.mu.Lock()
	r.series = nil
	r.mu.Unlock()

	go r.run()
	return nil
}

// Stop 请求停止（不阻塞；由 run 收尾）
func (r *Runner) Stop() {
	if r.status.Load() != StatusRunning {
		return
	}
	r.status.CompareAndSwap(StatusRunning, StatusStopping)
	if r.cancel != nil {
		r.cancel()
	}
	if r.pub != nil {
		r.pub.Stop()
	}
}

// run 主流程（goroutine 内执行）
func (r *Runner) run() {
	defer func() {
		// 收尾
		r.pool.Close()
		r.stopAt = time.Now()
		r.status.Store(StatusIdle)
		r.emit(r.buildSnapshot())
	}()

	// 1. 连接 + 订阅放 goroutine：连接期间采样循环同步推送进度
	connDone := make(chan struct{})
	go func() {
		connOK := r.pool.Start(r.ctx)
		r.connOK = int64(connOK)
		close(connDone)
	}()
	// 发布器等连接完成后启动
	go func() {
		<-connDone
		if r.cfg.PublishEnabled && r.pub != nil && r.connOK > 0 {
			r.pub.Start(r.ctx)
		}
	}()

	// 2. 运行至：时长到 / 手动 Stop / ctx 取消
	durCh := make(chan struct{})
	if r.cfg.Duration > 0 {
		go func() {
			select {
			case <-time.After(time.Duration(r.cfg.Duration) * time.Second):
				close(durCh)
			case <-r.ctx.Done():
			}
		}()
	} else {
		close(durCh)
	}

	// 3. 采样推送循环（连接阶段起即推送，连接进度实时可见）
	sampleTicker := time.NewTicker(500 * time.Millisecond)
	defer sampleTicker.Stop()
loop:
	for {
		select {
		case <-r.ctx.Done():
			break loop
		case <-durCh:
			break loop
		case <-sampleTicker.C:
			r.sample()
		}
	}
	r.Stop() // 触发收尾
}

// sample 采样一次：计算速率、追加时间序列、推送前端
func (r *Runner) sample() {
	now := time.Now()
	dt := now.Sub(r.lastSampleTime).Seconds()
	if dt <= 0 {
		dt = 0.5
	}
	pubTotal := r.metrics.PublishTotal.Load()
	recvTotal := r.metrics.ReceivedTotal.Load()
	pubRate := int64(float64(pubTotal-r.lastPub) / dt)
	recvRate := int64(float64(recvTotal-r.lastRecv) / dt)
	r.lastPub, r.lastRecv = pubTotal, recvTotal
	r.lastSampleTime = now

	raw := r.metrics.buildRaw(r.statusText(), r.pool.ConnectedCount())
	point := map[string]any{
		"ts":          now.UnixMilli(),
		"connected":   raw.Connected,
		"pub_rate":    pubRate,
		"recv_rate":   recvRate,
		"echo_avg_ms": raw.EchoAvgMs,
	}
	r.mu.Lock()
	r.series = append(r.series, point)
	if len(r.series) > 1200 { // 最多保留 10 分钟
		r.series = r.series[len(r.series)-1200:]
	}
	r.mu.Unlock()

	r.emit(r.buildSnapshot())
}

// statusText 状态文案
func (r *Runner) statusText() string {
	switch r.Status() {
	case StatusRunning:
		return "running"
	case StatusStopping:
		return "stopping"
	default:
		return "idle"
	}
}

// buildSnapshot 组装完整快照 JSON 字符串（推送给前端）
func (r *Runner) buildSnapshot() string {
	raw := r.metrics.buildRaw(r.statusText(), r.pool.ConnectedCount())
	raw.StartAt = r.startAt.UnixMilli()
	raw.DurationSec = int64(r.cfg.Duration)
	r.mu.Lock()
	raw.Series = append([]map[string]any(nil), r.series...)
	r.mu.Unlock()
	b, err := json.Marshal(raw)
	if err != nil {
		return fmt.Sprintf(`{"status":"%s","error":"%v"}`, r.statusText(), err)
	}
	return string(b)
}

// SnapshotMap 供 GetStatus 返回（直接 map）
func (r *Runner) SnapshotMap() map[string]any {
	raw := r.metrics.buildRaw(r.statusText(), r.pool.ConnectedCount())
	r.mu.Lock()
	raw.Series = append([]map[string]any(nil), r.series...)
	r.mu.Unlock()
	return map[string]any{
		"status":            raw.Status,
		"start_at":          r.startAt.Format("2006-01-02 15:04:05"),
		"start_at_ms":       r.startAt.UnixMilli(),
		"stop_at":           r.stopAt.Format("2006-01-02 15:04:05"),
		"duration_sec":      int(r.stopAt.Sub(r.startAt).Seconds()),
		"connected":         raw.Connected,
		"connect_total":     raw.ConnectTotal,
		"connect_success":   raw.ConnectSuccess,
		"connect_fail":      raw.ConnectFail,
		"connect_avg_ms":    raw.ConnectAvgMs,
		"connect_max_ms":    raw.ConnectMaxMs,
		"subscribe_total":   raw.SubscribeTotal,
		"subscribe_success": raw.SubscribeSucc,
		"subscribe_fail":    raw.SubscribeFail,
		"publish_total":     raw.PublishTotal,
		"publish_success":   raw.PublishSuccess,
		"publish_fail":      raw.PublishFail,
		"received_total":    raw.ReceivedTotal,
		"echo_count":        raw.EchoCount,
		"echo_avg_ms":       raw.EchoAvgMs,
		"echo_max_ms":       raw.EchoMaxMs,
		"echo_p50_ms":       raw.EchoP50Ms,
		"echo_p95_ms":       raw.EchoP95Ms,
		"echo_p99_ms":       raw.EchoP99Ms,
		"fail_reasons":      sortFailReasons(raw.FailReasons),
		"series":            raw.Series,
		"config":            r.cfg,
	}
}

// emit 推送事件（空函数兜底，避免 nil 调用）
func (r *Runner) emit(payload string) {
	if r.emitFn != nil {
		r.emitFn("stress:metrics", payload)
	}
}

// ExportCSV 导出结果到指定目录，返回文件路径
func (r *Runner) ExportCSV(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("stress_result_%s.csv", r.stopAt.Format("20060102_150405"))
	if r.stopAt.IsZero() {
		name = fmt.Sprintf("stress_result_%s.csv", time.Now().Format("20060102_150405"))
	}
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	raw := r.metrics.buildRaw("idle", r.pool.ConnectedCount())
	rows := [][]string{
		{"MQTTW 压测结果", name},
		{"开始时间", r.startAt.Format("2006-01-02 15:04:05")},
		{"结束时间", r.stopAt.Format("2006-01-02 15:04:05")},
		{"", ""},
		{"指标", "数值"},
		{"已连接", fmt.Sprintf("%d", raw.Connected)},
		{"连接总数(尝试)", fmt.Sprintf("%d", raw.ConnectTotal)},
		{"连接成功", fmt.Sprintf("%d", raw.ConnectSuccess)},
		{"连接失败", fmt.Sprintf("%d", raw.ConnectFail)},
		{"连接平均耗时ms", fmt.Sprintf("%d", raw.ConnectAvgMs)},
		{"连接最大耗时ms", fmt.Sprintf("%d", raw.ConnectMaxMs)},
		{"订阅总数", fmt.Sprintf("%d", raw.SubscribeTotal)},
		{"订阅成功", fmt.Sprintf("%d", raw.SubscribeSucc)},
		{"订阅失败", fmt.Sprintf("%d", raw.SubscribeFail)},
		{"发布总数", fmt.Sprintf("%d", raw.PublishTotal)},
		{"发布成功", fmt.Sprintf("%d", raw.PublishSuccess)},
		{"发布失败", fmt.Sprintf("%d", raw.PublishFail)},
		{"接收总数", fmt.Sprintf("%d", raw.ReceivedTotal)},
		{"回显次数", fmt.Sprintf("%d", raw.EchoCount)},
		{"回显平均延迟ms", fmt.Sprintf("%d", raw.EchoAvgMs)},
		{"回显最大延迟ms", fmt.Sprintf("%d", raw.EchoMaxMs)},
		{"回显P50延迟ms", fmt.Sprintf("%d", raw.EchoP50Ms)},
		{"回显P95延迟ms", fmt.Sprintf("%d", raw.EchoP95Ms)},
		{"回显P99延迟ms", fmt.Sprintf("%d", raw.EchoP99Ms)},
		{"", ""},
		{"失败原因", "次数"},
	}
	for k, v := range sortFailReasons(raw.FailReasons) {
		rows = append(rows, []string{fmt.Sprintf("%v", k), fmt.Sprintf("%v", v)})
	}
	rows = append(rows, []string{"", ""}, []string{"配置项", "值"})
	cfgB, _ := json.MarshalIndent(r.cfg, "", "  ")
	rows = append(rows, []string{"config", string(cfgB)})
	return path, w.WriteAll(rows)
}
