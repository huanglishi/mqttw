// Package modbus MQTTW Modbus→MQTT 网关采集核心（无 Wails/GUI 依赖，可独立打包为常驻服务）
package modbus

// 运行状态
const (
	StatusIdle     = 0
	StatusRunning  = 1
	StatusStopping = 2
)

// 寄存器类型（对应功能码 01/02/03/04）
const (
	RegisterCoil          = "coil"           // 0x 线圈
	RegisterDiscreteInput = "discrete_input" // 1x 离散输入
	RegisterHolding       = "holding"        // 3x 保持寄存器
	RegisterInput         = "input"          // 4x 输入寄存器
)

// 点位数据类型
const (
	DataTypeUint16  = "uint16"
	DataTypeInt16   = "int16"
	DataTypeUint32  = "uint32"
	DataTypeInt32   = "int32"
	DataTypeFloat32 = "float32"
	DataTypeBCD     = "bcd"
	DataTypeBit     = "bit"
)

// 字节序（单个 word 内部）
const (
	ByteOrderBig    = "big"    // 高字节在前（Modbus 默认）
	ByteOrderLittle = "little" // 低字节在前
)

// 字序（多寄存器组合）
const (
	WordOrderHighFirst = "high_first" // 起始寄存器为高字（float32 常见 ABCD）
	WordOrderLowFirst  = "low_first"  // 起始寄存器为低字（CDAB）
)

// 设备连接协议
const (
	ProtocolTCP        = "tcp"
	ProtocolRTU        = "rtu"
	ProtocolRTUOverTCP = "rtu_over_tcp"
)

// DeviceConfig 设备采集配置（service 层由 modbus_device 表转换而来）
type DeviceConfig struct {
	ID                 int32
	Name               string
	Protocol           string
	Host               string
	Port               int32
	SlaveID            int32
	SerialPort         string
	BaudRate           int32
	DataBits           int32
	StopBits           int32
	Parity             string // none / even / odd
	PollInterval       int32  // 采集周期 ms
	Timeout            int32  // 单次请求超时 ms
	RetryCount         int32  // 失败重试次数
	ByteOrder          string
	WordOrder          string
	BrokerConnectionID int32 // 关联 mqtt_connection.id，0 表示使用全局网关连接
	TopicPrefix        string
	Qos                int32
	Retain             int32
	AlarmTopic         string
}

// PointConfig 点位映射配置（service 层由 modbus_point 表转换而来）
type PointConfig struct {
	ID            int32
	DeviceID      int32
	Name          string
	RegisterType  string
	Address       int32
	Quantity      int32
	DataType      string
	Scale         float64
	Offset        float64
	Unit          string
	TopicOverride string // 非空时该点位独立主题发布
	OnlyOnChange  int32
	Enabled       int32
	Sort          int32
}

// PointValue 单点位读取结果（测试读取 / 实时面板用）
type PointValue struct {
	PointID   int32   `json:"point_id"`
	Name      string  `json:"name"`
	Register  string  `json:"register_type"`
	Address   int32   `json:"address"`
	DataType  string  `json:"data_type"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit,omitempty"`
	Raw       string  `json:"raw,omitempty"` // 原始字节 HEX
	Timestamp int64   `json:"ts"`
}

// DeviceStatus 设备运行状态（实时面板）
type DeviceStatus struct {
	DeviceID       int32   `json:"device_id"`
	Name           string  `json:"name"`
	Protocol       string  `json:"protocol"`
	Running        bool    `json:"running"`
	Connected      bool    `json:"connected"`
	LastError      string  `json:"last_error,omitempty"`
	PollCount      int64   `json:"poll_count"`
	ErrorCount     int64   `json:"error_count"`
	LastPollAt     int64   `json:"last_poll_at"`
	PollInterval   int32   `json:"poll_interval"`
	ValuesPerCycle int     `json:"values_per_cycle"`
	Rate           float64 `json:"rate"` // 最近 1 分钟平均每秒采集点数
}
