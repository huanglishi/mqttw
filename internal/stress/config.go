// 压测配置：前端表单 JSON 直接映射
package stress

import (
	"encoding/json"
	"fmt"
	"time"
)

// StressConfig 压测任务配置（与前端 StressConfig.vue 表单字段一一对应）
type StressConfig struct {
	// 连接
	Broker         string `json:"broker"`     // tcp://host:port
	Protocol       string `json:"protocol"`   // 3.1 / 3.1.1 / 5.0 / auto
	Username       string `json:"username"`   //
	Password       string `json:"password"`   //
	ClientCount    int    `json:"clientCount"` // 总客户端数
	ClientPrefix   string `json:"clientPrefix"`
	ConnectRate    int    `json:"connectRate"`    // 每秒新建连接数，0=不限速
	ConnectTimeout int    `json:"connectTimeout"` // 秒
	KeepAlive      int    `json:"keepalive"`      // 秒
	CleanSession   bool   `json:"cleanSession"`

	// 订阅
	Topics             []string `json:"topics"` // 订阅模板（支持通配符）
	SubscribePerClient int      `json:"subscribePerClient"`
	Qos                int      `json:"qos"`

	// 发布
	PublishEnabled     bool   `json:"publishEnabled"`
	PublishWorkers     int    `json:"publishWorkers"`
	PublishTopic       string `json:"publishTopic"`
	PublishPayloadSize int    `json:"publishPayloadSize"` // 字节
	PublishQos         int    `json:"publishQos"`
	PublishRate        int    `json:"publishRate"`  // msg/s，0=不限
	PublishTotal       int64  `json:"publishTotal"` // 总条数，0=不限
	Duration           int    `json:"duration"`     // 秒，0=不限

	// 回显（端到端延迟测量）
	EchoEnabled      bool   `json:"echoEnabled"`
	PublishEchoTopic string `json:"publishEchoTopic"` // 回显订阅 topic，默认=PublishTopic
}

// normalize 补默认值并校验
func (c *StressConfig) normalize() error {
	if c.Broker == "" {
		return fmt.Errorf("broker 地址不能为空")
	}
	if c.ClientCount <= 0 {
		c.ClientCount = 1
	}
	if c.ClientPrefix == "" {
		c.ClientPrefix = "stress_test"
	}
	if c.ConnectTimeout <= 0 {
		c.ConnectTimeout = 5
	}
	if c.KeepAlive <= 0 {
		c.KeepAlive = 30
	}
	if c.ConnectRate < 0 {
		c.ConnectRate = 0
	}
	switch c.Protocol {
	case "", "auto":
		c.Protocol = "auto"
	case "3.1", "3.1.1", "5.0":
	default:
		return fmt.Errorf("不支持的协议版本: %s（仅支持 3.1 / 3.1.1 / 5.0 / auto）", c.Protocol)
	}
	if c.SubscribePerClient < 0 {
		c.SubscribePerClient = 0
	}
	if len(c.Topics) > 0 && c.SubscribePerClient > len(c.Topics) {
		c.SubscribePerClient = len(c.Topics)
	}
	if c.PublishWorkers <= 0 {
		c.PublishWorkers = 4
	}
	if c.PublishPayloadSize < 0 {
		c.PublishPayloadSize = 0
	}
	if c.PublishRate < 0 {
		c.PublishRate = 0
	}
	if c.PublishTotal < 0 {
		c.PublishTotal = 0
	}
	if c.Duration < 0 {
		c.Duration = 0
	}
	// 默认 topic
	if len(c.Topics) == 0 {
		c.Topics = []string{"stress/#"}
	}
	if c.PublishTopic == "" {
		c.PublishTopic = "stress/pub"
	}
	if c.PublishEchoTopic == "" {
		c.PublishEchoTopic = c.PublishTopic
	}
	return nil
}

// ParseConfig 解析前端 JSON 并归一化
func ParseConfig(param string) (StressConfig, error) {
	var cfg StressConfig
	if err := json.Unmarshal([]byte(param), &cfg); err != nil {
		return cfg, fmt.Errorf("压测配置解析失败: %w", err)
	}
	if err := cfg.normalize(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// connectTimeoutDuration 连接超时
func (c *StressConfig) connectTimeoutDuration() time.Duration {
	return time.Duration(c.ConnectTimeout) * time.Second
}

// connectInterval 连接限速间隔；0=不限速
func (c *StressConfig) connectInterval() time.Duration {
	if c.ConnectRate <= 0 {
		return 0
	}
	return time.Second / time.Duration(c.ConnectRate)
}
