// 抽象、结构体、工厂方法
package mqttclient

import "context"

type ProtocolVersion string

const (
	ProtocolV31  ProtocolVersion = "3.1"
	ProtocolV311 ProtocolVersion = "3.1.1"
	ProtocolV50  ProtocolVersion = "5.0"
)

// ClientConfig 客户端配置
type ClientConfig struct {
	Broker        string
	ClientID      string
	Username      string
	Password      string
	Protocol      ProtocolVersion
	KeepAlive     uint16
	WillTopic     string
	WillPayload   []byte
	WillQos       byte
	WillRetain    bool
	TLSEnable     bool
	CleanSession  bool // MQTT3.x
	CleanStart    bool // MQTT5
	SessionExpiry uint32

	AutoReconnect   bool
	ReconnectPeriod uint16

	// MQTT5 Connect 属性，仅5.0生效
	ReceiveMaximum             uint16
	MaximumPacketSize          uint32
	TopicAliasMaximum          uint16
	RequestResponseInformation bool
	RequestProblemInformation  bool
	UserProperties             [][]string
}

// IMqttClient 统一调试客户端接口
type IMqttClient interface {
	Connect(ctx context.Context) error
	Disconnect() error
	Subscribe(ctx context.Context, opt SubscribeOption) error
	UnSubscribe(ctx context.Context, topic string) error
	Publish(ctx context.Context, opt PublishOption) error
	OnMessage(cb func(msg *MqttMessage))
	IsConnected() bool
}

// SubscribeOption 订阅参数
type SubscribeOption struct {
	Topic                  string
	Qos                    byte
	SubscriptionIdentifier uint32
	NoLocal                bool
	RetainAsPublished      bool
	RetainHandling         byte
}

// PublishOption 发布消息参数
type PublishOption struct {
	Topic                  string
	Payload                []byte
	Qos                    byte
	Retain                 bool
	UserProperties         [][]string
	ResponseTopic          string
	ContentType            string
	CorrelationData        []byte
	MessageExpiry          uint32
	TopicAlias             uint16
	PayloadFormat          byte
	SubscriptionIdentifier int
}

// MqttMessage 统一对外消息结构体，给到Wails事件
type MqttMessage struct {
	Topic                   string
	Payload                 []byte
	Qos                     byte
	Retain                  bool
	SubscriptionIdentifiers []uint32
	UserProperties          [][]string
	ResponseTopic           string
	ContentType             string
	CorrelationData         []byte
	MessageExpiry           uint32
}

// NewMqttClient 调试客户端工厂
func NewMqttClient(cfg ClientConfig) IMqttClient {
	// MQTT3.x 清空MQTT5专属参数，避免干扰底层
	if cfg.Protocol != ProtocolV50 {
		cfg.ReceiveMaximum = 0
		cfg.MaximumPacketSize = 0
		cfg.TopicAliasMaximum = 0
		cfg.RequestResponseInformation = false
		cfg.RequestProblemInformation = false
		cfg.UserProperties = nil
	}

	switch cfg.Protocol {
	case ProtocolV31, ProtocolV311:
		return NewPahoClient(cfg)
	case ProtocolV50:
		return NewPahoV5Client(cfg)
	default:
		return nil
	}
}
