package modbus

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"gofly/internal/mqttclient"
)

// Gateway 网关编排器：管理多台设备的采集协程与 MQTT 输出。
// Publisher 按 broker_connection_id 分桶共享：相同 Broker 的设备复用同一连接。
type Gateway struct {
	mu         sync.Mutex
	collectors map[int32]*Collector
	pubs       map[int32]*Publisher // broker_connection_id -> 发布器（0 表示默认/未配置）

	emit    func(event, payload string)
	status  atomic.Int32
	startAt time.Time
}

// NewGateway 创建网关编排器
func NewGateway(emit func(event, payload string)) *Gateway {
	return &Gateway{
		collectors: make(map[int32]*Collector),
		pubs:       make(map[int32]*Publisher),
		emit:       emit,
	}
}

// Start 启动网关：为每个启用设备创建采集器并启动；按 Broker 桶创建 MQTT 发布器。
// devices: 设备配置；pointsByDevice: device_id -> 点位列表；
// mqttByBroker: broker_connection_id -> MQTT 客户端配置（缺失则设备仅采集不上报）。
func (g *Gateway) Start(devices []DeviceConfig, pointsByDevice map[int32][]PointConfig, mqttByBroker map[int32]mqttclient.ClientConfig) error {
	g.mu.Lock()
	if g.status.Load() == StatusRunning {
		g.mu.Unlock()
		return fmt.Errorf("网关已在运行，请先停止")
	}
	g.status.Store(StatusRunning)
	g.startAt = time.Now()
	g.mu.Unlock()

	ctx := context.Background()
	started := 0
	for _, cfg := range devices {
		points := pointsByDevice[cfg.ID]
		var pub *Publisher
		if mqttCfg, ok := mqttByBroker[cfg.BrokerConnectionID]; ok {
			pub = g.publisher(cfg.BrokerConnectionID, mqttCfg, ctx)
		}
		col := NewCollector(cfg, points, pub, g.emit)
		if err := col.Start(); err != nil {
			g.log("error", fmt.Sprintf("设备 %s 启动失败: %v", cfg.Name, err))
			continue
		}
		g.mu.Lock()
		g.collectors[cfg.ID] = col
		g.mu.Unlock()
		started++
	}

	if started == 0 {
		g.log("error", "没有可启动的设备，请检查设备配置与 enabled 状态")
	}
	g.log("info", fmt.Sprintf("网关启动完成：%d/%d 台设备进入采集", started, len(devices)))
	g.emitStatus()
	return nil
}

// publisher 获取或创建指定 Broker 桶的发布器
func (g *Gateway) publisher(brokerID int32, mqttCfg mqttclient.ClientConfig, ctx context.Context) *Publisher {
	g.mu.Lock()
	defer g.mu.Unlock()
	if p, ok := g.pubs[brokerID]; ok {
		return p
	}
	p := NewPublisher(mqttCfg)
	if err := p.Start(ctx); err != nil {
		g.log("warn", fmt.Sprintf("MQTT 连接失败（broker=%s）: %v，该组设备仅采集不上报", mqttCfg.Broker, err))
	} else {
		g.log("info", fmt.Sprintf("MQTT 输出已连接 %s", mqttCfg.Broker))
	}
	g.pubs[brokerID] = p
	return p
}

// Stop 停止全部采集与 MQTT 输出
func (g *Gateway) Stop() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.status.Load() != StatusRunning {
		return
	}
	g.status.Store(StatusStopping)
	for id, col := range g.collectors {
		col.Stop()
		delete(g.collectors, id)
	}
	for _, p := range g.pubs {
		p.Close()
	}
	g.pubs = map[int32]*Publisher{}
	g.status.Store(StatusIdle)
	g.log("info", "网关已停止")
	g.emitStatus()
}

// Running 是否运行中
func (g *Gateway) Running() bool { return g.status.Load() == StatusRunning }

// Status 网关状态快照
func (g *Gateway) Status() map[string]any {
	g.mu.Lock()
	defer g.mu.Unlock()
	devices := make([]DeviceStatus, 0, len(g.collectors))
	var totalPolls, totalErrors, totalPoints int64
	for id, col := range g.collectors {
		st := col.Status()
		st.DeviceID = id
		devices = append(devices, st)
		totalPolls += st.PollCount
		totalErrors += st.ErrorCount
		totalPoints += int64(st.ValuesPerCycle)
	}
	pubs := make([]map[string]any, 0, len(g.pubs))
	for _, p := range g.pubs {
		pubs = append(pubs, p.Status())
	}
	return map[string]any{
		"running":      g.status.Load() == StatusRunning,
		"start_at":     g.startAt.Unix(),
		"device_count": len(devices),
		"total_points": totalPoints,
		"total_polls":  totalPolls,
		"total_errors": totalErrors,
		"devices":      devices,
		"mqtt":         pubs,
	}
}

// TestRead 单次在线读取测试（不启动采集任务）：连接设备并读取所有点位，返回解析结果。
func (g *Gateway) TestRead(cfg DeviceConfig, points []PointConfig) ([]PointValue, error) {
	dev := NewDevice(cfg)
	if err := dev.Connect(); err != nil {
		return nil, err
	}
	defer dev.Close()

	pts := make([]PointConfig, 0, len(points))
	for _, p := range points {
		if p.Enabled == 0 {
			continue
		}
		pts = append(pts, p)
	}
	if len(pts) == 0 {
		return nil, fmt.Errorf("设备 %s 没有启用的点位可读取", cfg.Name)
	}

	out := make([]PointValue, 0, len(pts))
	for _, p := range pts {
		q := uint16(p.Quantity)
		if q <= 0 {
			q = 1
		}
		raw, err := dev.Read(p.RegisterType, uint16(p.Address), q)
		if err != nil {
			out = append(out, PointValue{
				PointID: p.ID, Name: p.Name, Register: p.RegisterType,
				Address: p.Address, DataType: p.DataType, Unit: p.Unit,
				Timestamp: time.Now().Unix(),
			})
			// 单点位失败不中断测试，最后一条携带错误
			g.log("warn", fmt.Sprintf("测试读取 设备 %s 点位 %s 失败: %v", cfg.Name, p.Name, err))
			continue
		}
		chunk := raw
		if p.RegisterType != RegisterCoil && p.RegisterType != RegisterDiscreteInput {
			if len(raw) < int(q)*2 {
				chunk = raw
			} else {
				chunk = raw[:int(q)*2]
			}
		} else {
			chunk = raw[:1]
		}
		v, err := ParseRegisters(chunk, p, cfg.ByteOrder, cfg.WordOrder)
		if err != nil {
			g.log("warn", fmt.Sprintf("测试读取 设备 %s 点位 %s 解析失败: %v", cfg.Name, p.Name, err))
			continue
		}
		out = append(out, PointValue{
			PointID:   p.ID,
			Name:      p.Name,
			Register:  p.RegisterType,
			Address:   p.Address,
			DataType:  p.DataType,
			Value:     round6(ApplyScale(v, p)),
			Unit:      p.Unit,
			Raw:       hexStr(raw),
			Timestamp: time.Now().Unix(),
		})
	}
	return out, nil
}

// log 网关日志事件
func (g *Gateway) log(level, msg string) {
	if g.emit == nil {
		return
	}
	data, _ := json.Marshal(map[string]any{
		"level": level,
		"msg":   msg,
		"ts":    time.Now().Unix(),
	})
	g.emit("gateway:log", string(data))
}

// emitStatus 推送状态事件
func (g *Gateway) emitStatus() {
	if g.emit == nil {
		return
	}
	data, _ := json.Marshal(g.Status())
	g.emit("gateway:status", string(data))
}

// hexStr 原始字节 HEX 显示
func hexStr(b []byte) string {
	const hex = "0123456789ABCDEF"
	out := make([]byte, 0, len(b)*3)
	for i, v := range b {
		if i > 0 {
			out = append(out, ' ')
		}
		out = append(out, hex[v>>4], hex[v&0x0f])
	}
	return string(out)
}
