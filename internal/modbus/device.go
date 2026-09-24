package modbus

import (
	"fmt"
	"strings"
	"sync"
	"time"

	mb "github.com/goburrow/modbus"
)

// Device 单台 Modbus 设备（连接 + 读写 + 状态统计）。
// 注意：Modbus RTU 串口同一时刻只能有一个请求，采集循环为单协程串行，
// 天然满足串口并发约束，无需额外锁队列。
type Device struct {
	cfg  DeviceConfig
	cli  mb.Client
	conn *mb.TCPClientHandler // TCP / RTU over TCP

	mu        sync.Mutex
	connected bool
	lastErr   string

	pollCount  int64
	errorCount int64
	lastPollAt int64
}

// NewDevice 创建设备实例（不建立连接）
func NewDevice(cfg DeviceConfig) *Device {
	return &Device{cfg: cfg}
}

// Config 返回设备配置
func (d *Device) Config() DeviceConfig { return d.cfg }

// Name 设备名
func (d *Device) Name() string { return d.cfg.Name }

// ID 设备 ID
func (d *Device) ID() int32 { return d.cfg.ID }

// Connect 建立底层连接（TCP 握手 / 串口打开）。
// RTU over TCP 本质是 TCP 连接 + SlaveId 作为单元标识（网关透传场景）。
func (d *Device) Connect() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.connected {
		return nil
	}

	switch d.cfg.Protocol {
	case ProtocolTCP, ProtocolRTUOverTCP:
		if d.cfg.Host == "" {
			return fmt.Errorf("TCP 设备缺少 host 地址")
		}
		addr := fmt.Sprintf("%s:%d", d.cfg.Host, d.cfg.Port)
		h := mb.NewTCPClientHandler(addr)
		h.Timeout = timeoutOf(d.cfg)
		h.IdleTimeout = 60 * time.Second
		h.SlaveId = slaveID(d.cfg)
		d.conn = h
		d.cli = mb.NewClient(h)
	case ProtocolRTU:
		if d.cfg.SerialPort == "" {
			return fmt.Errorf("RTU 设备缺少串口（如 COM3）")
		}
		h := mb.NewRTUClientHandler(d.cfg.SerialPort)
		h.Timeout = timeoutOf(d.cfg)
		h.IdleTimeout = 60 * time.Second
		h.SlaveId = slaveID(d.cfg)
		h.BaudRate = int(d.cfg.BaudRate)
		h.DataBits = int(d.cfg.DataBits)
		h.StopBits = int(d.cfg.StopBits)
		h.Parity = parityOf(d.cfg.Parity)
		if err := h.Connect(); err != nil {
			return fmt.Errorf("打开串口 %s 失败: %w", d.cfg.SerialPort, err)
		}
		d.cli = mb.NewClient(h)
	default:
		return fmt.Errorf("不支持的协议 %q", d.cfg.Protocol)
	}

	d.connected = true
	d.lastErr = ""
	return nil
}

// Close 关闭连接
func (d *Device) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.conn != nil {
		_ = d.conn.Close()
		d.conn = nil
	}
	d.cli = nil
	d.connected = false
}

// Connected 是否已连接
func (d *Device) Connected() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.connected
}

// Read 按寄存器类型读取 quantity 个寄存器/线圈，返回原始字节。
// 读失败时自动标记断连（下一次采集会触发重连）。
func (d *Device) Read(registerType string, address, quantity uint16) ([]byte, error) {
	if !d.Connected() {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	var (
		raw []byte
		err error
	)
	switch registerType {
	case RegisterCoil:
		raw, err = d.cli.ReadCoils(address, quantity)
	case RegisterDiscreteInput:
		raw, err = d.cli.ReadDiscreteInputs(address, quantity)
	case RegisterHolding:
		raw, err = d.cli.ReadHoldingRegisters(address, quantity)
	case RegisterInput:
		raw, err = d.cli.ReadInputRegisters(address, quantity)
	default:
		err = fmt.Errorf("不支持的寄存器类型 %q", registerType)
	}

	d.mu.Lock()
	if err != nil {
		d.errorCount++
		d.lastErr = err.Error()
		// 网络/串口类错误视为断连，下次采集前自动重连
		if isTransientError(err) {
			d.connected = false
		}
	} else {
		d.pollCount++
		d.lastPollAt = time.Now().Unix()
		d.lastErr = ""
	}
	d.mu.Unlock()
	return raw, err
}

// Status 设备状态快照
func (d *Device) Status() DeviceStatus {
	d.mu.Lock()
	defer d.mu.Unlock()
	return DeviceStatus{
		DeviceID:     d.cfg.ID,
		Name:         d.cfg.Name,
		Protocol:     d.cfg.Protocol,
		Connected:    d.connected,
		LastError:    d.lastErr,
		PollCount:    d.pollCount,
		ErrorCount:   d.errorCount,
		LastPollAt:   d.lastPollAt,
		PollInterval: d.cfg.PollInterval,
	}
}

func timeoutOf(cfg DeviceConfig) time.Duration {
	t := time.Duration(cfg.Timeout) * time.Millisecond
	if t <= 0 {
		t = 1 * time.Second
	}
	return t
}

func slaveID(cfg DeviceConfig) byte {
	id := cfg.SlaveID
	if id <= 0 {
		id = 1
	}
	return byte(id)
}

// parityOf 归一化串口校验位为 goburrow/serial 期望的字符（N/E/O）
func parityOf(p string) string {
	switch strings.ToLower(strings.TrimSpace(p)) {
	case "even":
		return "E"
	case "odd":
		return "O"
	default:
		return "N"
	}
}

// isTransientError 判断错误是否属于可重连的瞬时错误
func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, kw := range []string{
		"timeout", "timed out", "connection reset", "broken pipe",
		"connection refused", "i/o timeout", "port is closed",
		"no such host", "network is unreachable", "eof",
	} {
		if strings.Contains(msg, kw) {
			return true
		}
	}
	return false
}
