package modbus

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// 批量读上限（Modbus 协议单请求最大 125 个寄存器 / 2000 个线圈）
const (
	maxRegistersPerBatch = 125
	maxBitsPerBatch      = 2000
)

// Collector 单台设备的采集任务（独立协程，串行轮询）。
// RTU 串口同一时刻只能一个请求：单协程串行天然满足，无需串口锁队列。
type Collector struct {
	cfg    DeviceConfig
	dev    *Device
	points []PointConfig // 已过滤 enabled 并按 sort 排序
	pub    *Publisher
	emit   func(event, payload string)

	status atomic.Int32
	ctx    context.Context
	cancel context.CancelFunc

	mu       sync.Mutex
	last     map[string]float64 // 变化检测缓存：name -> 上次值
	values   map[string]float64 // 最新一轮全量值（面板）
	cyclePts int

	rateMu     sync.Mutex
	rateWin    [10]float64
	rateIdx    int
	startAt    time.Time
	lastCycle  time.Time
	lastAlarm  time.Time
	alarmCount int64
}

// NewCollector 创建采集器。pub 可为 nil（仅测试读取场景）。
func NewCollector(cfg DeviceConfig, points []PointConfig, pub *Publisher, emit func(event, payload string)) *Collector {
	enabled := make([]PointConfig, 0, len(points))
	for _, p := range points {
		if p.Enabled == 0 {
			continue
		}
		enabled = append(enabled, p)
	}
	sort.SliceStable(enabled, func(i, j int) bool {
		if enabled[i].RegisterType != enabled[j].RegisterType {
			return enabled[i].RegisterType < enabled[j].RegisterType
		}
		if enabled[i].Address != enabled[j].Address {
			return enabled[i].Address < enabled[j].Address
		}
		return enabled[i].Sort < enabled[j].Sort
	})
	return &Collector{
		cfg:       cfg,
		dev:       NewDevice(cfg),
		points:    enabled,
		pub:       pub,
		emit:      emit,
		last:      make(map[string]float64),
		values:    make(map[string]float64),
		startAt:   time.Now(),
		lastCycle: time.Now(),
	}
}

// Start 启动采集协程（非阻塞）
func (c *Collector) Start() error {
	if !c.status.CompareAndSwap(StatusIdle, StatusRunning) {
		return fmt.Errorf("设备 %s 已在运行", c.cfg.Name)
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	interval := time.Duration(c.cfg.PollInterval) * time.Millisecond
	if interval <= 0 {
		interval = time.Second
	}
	go c.loop(interval)
	return nil
}

// Stop 停止采集
func (c *Collector) Stop() {
	if c.status.Load() != StatusRunning {
		return
	}
	c.status.CompareAndSwap(StatusRunning, StatusStopping)
	if c.cancel != nil {
		c.cancel()
	}
	// 等待退出后复位状态由 loop defer 完成
}

// loop 轮询主循环
func (c *Collector) loop(interval time.Duration) {
	defer func() {
		c.status.Store(StatusIdle)
		c.dev.Close()
		c.log("info", fmt.Sprintf("设备 %s 采集已停止", c.cfg.Name))
	}()

	if err := c.dev.Connect(); err != nil {
		c.log("error", fmt.Sprintf("设备 %s 连接失败: %v", c.cfg.Name, err))
	} else {
		c.log("info", fmt.Sprintf("设备 %s 已连接，开始采集（周期 %dms，点位 %d）", c.cfg.Name, c.cfg.PollInterval, len(c.points)))
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.pollOnce()
		}
	}
}

// pollOnce 单轮采集：分组批量读 → 解析 → 变化检测 → 发布
func (c *Collector) pollOnce() {
	start := time.Now()
	values := make(map[string]float64)
	devValues := make(map[string]any) // 设备级 payload（无 topic_override 的点）
	changed := false

	for _, g := range groupPoints(c.points) {
		batchValues, ok := c.readBatch(g)
		if !ok {
			continue // 该组读取失败，跳过（已记日志/告警）
		}
		for name, v := range batchValues {
			values[name] = v
			p := c.pointByName(name)
			if p == nil {
				continue
			}
			// 变化检测：only_on_change 且值未变 → 跳过发布
			if p.OnlyOnChange == 1 {
				c.mu.Lock()
				last, exists := c.last[name]
				c.mu.Unlock()
				if exists && sameValue(last, v) {
					continue
				}
				c.mu.Lock()
				c.last[name] = v
				c.mu.Unlock()
				changed = true
			}
			// topic_override：独立主题发布单点
			if p.TopicOverride != "" {
				_ = c.publishPoint(*p, v)
				continue
			}
			devValues[name] = v
		}
	}

	// 设备级 payload 发布
	if c.pub != nil && len(devValues) > 0 {
		payload := map[string]any{"deviceId": c.cfg.Name}
		for k, v := range devValues {
			payload[k] = v
		}
		payload["ts"] = time.Now().Unix()
		if err := c.publishDevice(payload); err != nil {
			c.log("warn", fmt.Sprintf("设备 %s MQTT 发布失败: %v", c.cfg.Name, err))
		}
	}
	_ = changed

	// 面板事件：全量最新值
	c.mu.Lock()
	c.values = values
	c.cyclePts = len(values)
	c.mu.Unlock()
	c.updateRate(len(values), time.Since(start))
	if c.emit != nil {
		data, _ := json.Marshal(map[string]any{
			"device_id": c.cfg.ID,
			"device":    c.cfg.Name,
			"values":    values,
			"ts":        time.Now().Unix(),
		})
		c.emit("gateway:values", string(data))
	}
}

// readBatch 读取一组（同寄存器类型、连续地址合并）点位，返回 name->解析值。
// 失败时按 retry_count 重试并记录告警。
func (c *Collector) readBatch(g pointGroup) (map[string]float64, bool) {
	out := make(map[string]float64)
	for _, b := range g.batches {
		raw, err := c.readWithRetry(b)
		if err != nil {
			c.log("error", fmt.Sprintf("设备 %s 读取 %s 地址 %d 失败: %v", c.cfg.Name, g.regType, b.start, err))
			c.alarmIfThrottled(fmt.Sprintf("设备 %s 读取失败（%s@%d）: %v", c.cfg.Name, g.regType, b.start, err))
			return nil, false
		}
		// 切片解析批内点位
		for i, p := range b.points {
			chunk, err := sliceChunk(raw, b, p, i)
			if err != nil {
				c.log("warn", fmt.Sprintf("设备 %s 点位 %s 解析切片失败: %v", c.cfg.Name, p.Name, err))
				continue
			}
			v, err := ParseRegisters(chunk, p, c.cfg.ByteOrder, c.cfg.WordOrder)
			if err != nil {
				c.log("warn", fmt.Sprintf("设备 %s 点位 %s 解析失败: %v", c.cfg.Name, p.Name, err))
				continue
			}
			v = round6(ApplyScale(v, p))
			out[p.Name] = v
		}
	}
	return out, true
}

// readWithRetry 带重试的批量读
func (c *Collector) readWithRetry(b batch) ([]byte, error) {
	retries := int(c.cfg.RetryCount)
	if retries <= 0 {
		retries = 1
	}
	var lastErr error
	for i := 0; i < retries; i++ {
		if c.ctx.Err() != nil {
			return nil, c.ctx.Err()
		}
		raw, err := c.dev.Read(b.regType, uint16(b.start), uint16(b.quantity))
		if err == nil {
			return raw, nil
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	return nil, lastErr
}

// publishDevice 发布设备级 payload
func (c *Collector) publishDevice(payload map[string]any) error {
	if c.pub == nil {
		return nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.pub.Publish(ctx, c.deviceTopic(), data, byte(c.cfg.Qos), c.cfg.Retain == 1)
}

// publishPoint 发布独立主题点位
func (c *Collector) publishPoint(p PointConfig, v float64) error {
	if c.pub == nil {
		return nil
	}
	payload, _ := json.Marshal(map[string]any{
		"deviceId": c.cfg.Name,
		"name":     p.Name,
		"value":    v,
		"unit":     p.Unit,
		"ts":       time.Now().Unix(),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.pub.Publish(ctx, p.TopicOverride, payload, byte(c.cfg.Qos), c.cfg.Retain == 1)
}

// deviceTopic 设备默认主题：prefix 为空时使用 modbus/<设备名>
func (c *Collector) deviceTopic() string {
	if c.cfg.TopicPrefix != "" {
		return c.cfg.TopicPrefix
	}
	return "modbus/" + c.cfg.Name
}

func (c *Collector) pointByName(name string) *PointConfig {
	for i := range c.points {
		if c.points[i].Name == name {
			return &c.points[i]
		}
	}
	return nil
}

func (c *Collector) alarmIfThrottled(msg string) {
	now := time.Now()
	c.mu.Lock()
	throttled := now.Sub(c.lastAlarm) < 10*time.Second
	c.lastAlarm = now
	c.alarmCount++
	c.mu.Unlock()
	if throttled {
		return
	}
	c.log("alarm", msg)
	if c.pub != nil && c.cfg.AlarmTopic != "" {
		payload, _ := json.Marshal(map[string]any{
			"device": c.cfg.Name,
			"level":  "alarm",
			"msg":    msg,
			"ts":     time.Now().Unix(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = c.pub.Publish(ctx, c.cfg.AlarmTopic, payload, 0, false)
	}
}

// updateRate 滑动窗口速率（最近 10 轮实际每秒点数）
func (c *Collector) updateRate(n int, elapsed time.Duration) {
	sec := elapsed.Seconds()
	if sec <= 0 {
		sec = 0.001
	}
	c.rateMu.Lock()
	c.rateWin[c.rateIdx%len(c.rateWin)] = float64(n) / sec
	c.rateIdx++
	c.rateMu.Unlock()
}

// Rate 最近平均速率
func (c *Collector) Rate() float64 {
	c.rateMu.Lock()
	defer c.rateMu.Unlock()
	var sum float64
	n := c.rateIdx
	if n > len(c.rateWin) {
		n = len(c.rateWin)
	}
	if n == 0 {
		return 0
	}
	for i := 0; i < n; i++ {
		sum += c.rateWin[i]
	}
	return sum / float64(n)
}

// Status 采集器运行状态
func (c *Collector) Status() DeviceStatus {
	st := c.dev.Status()
	st.Running = c.status.Load() == StatusRunning
	st.Rate = c.Rate()
	c.mu.Lock()
	st.ValuesPerCycle = c.cyclePts
	c.mu.Unlock()
	return st
}

// Values 最新一轮全量值（实时面板）
func (c *Collector) Values() map[string]float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]float64, len(c.values))
	for k, v := range c.values {
		out[k] = v
	}
	return out
}

// log 采集日志事件
func (c *Collector) log(level, msg string) {
	if c.emit == nil {
		return
	}
	data, _ := json.Marshal(map[string]any{
		"device_id": c.cfg.ID,
		"device":    c.cfg.Name,
		"level":     level,
		"msg":       msg,
		"ts":        time.Now().Unix(),
	})
	c.emit("gateway:log", string(data))
}

// ---------- 点位分组与批量规划 ----------

type batch struct {
	regType  string // 寄存器类型
	start    uint16 // 起始地址
	quantity uint16 // 本批数量（寄存器数或线圈数）
	points   []PointConfig
}

type pointGroup struct {
	regType string
	batches []batch
}

// groupPoints 按寄存器类型分组，组内按地址排序后合并连续地址为批量读。
// 寄存器类合并上限 125，线圈类 2000；寄存器点位按各自 quantity 累加。
func groupPoints(points []PointConfig) []pointGroup {
	byType := map[string][]PointConfig{}
	var order []string
	for _, p := range points {
		if _, ok := byType[p.RegisterType]; !ok {
			order = append(order, p.RegisterType)
		}
		byType[p.RegisterType] = append(byType[p.RegisterType], p)
	}

	var groups []pointGroup
	for _, rt := range order {
		pts := byType[rt]
		sort.SliceStable(pts, func(i, j int) bool {
			if pts[i].Address != pts[j].Address {
				return pts[i].Address < pts[j].Address
			}
			return pts[i].Sort < pts[j].Sort
		})
		g := pointGroup{regType: rt}
		cur := batch{}
		limit := uint16(maxRegistersPerBatch)
		if rt == RegisterCoil || rt == RegisterDiscreteInput {
			limit = uint16(maxBitsPerBatch)
		}
		for i, p := range pts {
			q := uint16(1)
			if rt != RegisterCoil && rt != RegisterDiscreteInput {
				q = uint16(p.Quantity)
				if q <= 0 {
					q = 1
				}
			}
			if i == 0 {
				cur = batch{regType: rt, start: uint16(p.Address), quantity: q, points: []PointConfig{p}}
				continue
			}
			prev := pts[i-1]
			prevQ := uint16(1)
			if rt != RegisterCoil && rt != RegisterDiscreteInput {
				prevQ = uint16(prev.Quantity)
				if prevQ <= 0 {
					prevQ = 1
				}
			}
			nextStart := uint16(prev.Address) + prevQ
			merged := uint16(p.Address) >= cur.start && uint16(p.Address) <= nextStart // 连续或重叠
			if merged && cur.quantity+q <= limit {
				cur.quantity = maxU16(nextStart+q, cur.quantity)
				cur.points = append(cur.points, p)
				continue
			}
			g.batches = append(g.batches, cur)
			cur = batch{regType: rt, start: uint16(p.Address), quantity: q, points: []PointConfig{p}}
		}
		g.batches = append(g.batches, cur)
		groups = append(groups, g)
	}
	return groups
}

// sliceChunk 从批量读取结果中切出第 i 个点位的原始数据
func sliceChunk(raw []byte, b batch, p PointConfig, i int) ([]byte, error) {
	if b.regType == RegisterCoil || b.regType == RegisterDiscreteInput {
		// 位打包：每 8 个线圈 1 字节
		byteIdx := i / 8
		bitIdx := uint(i % 8)
		if byteIdx >= len(raw) {
			return nil, fmt.Errorf("线圈数据越界")
		}
		return []byte{(raw[byteIdx] >> bitIdx) & 0x01}, nil
	}
	// 寄存器：按批内偏移切 quantity*2 字节
	offset := 0
	for j := 0; j < i; j++ {
		q := int(b.points[j].Quantity)
		if q <= 0 {
			q = 1
		}
		offset += q * 2
	}
	end := offset + int(p.Quantity)*2
	if offset < 0 || end > len(raw) {
		return nil, fmt.Errorf("寄存器数据越界: offset=%d end=%d len=%d", offset, end, len(raw))
	}
	return raw[offset:end], nil
}

// ---------- 工具 ----------

func sameValue(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return a == b
}

// round6 四舍五入到 6 位小数，消除 float32→float64 精度噪声
func round6(v float64) float64 {
	return math.Round(v*1e6) / 1e6
}

func maxU16(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}
