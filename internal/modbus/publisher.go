package modbus

import (
	"context"
	"fmt"
	"sync"

	"gofly/internal/mqttclient"
)

// Publisher MQTT 输出层：为网关建立独立 MQTT 连接（不占用调试连接池），
// 复用 mqttclient 工厂（3.x / 5.0 双实现），支持遗嘱、离线告警。
type Publisher struct {
	cfg mqttclient.ClientConfig

	mu        sync.Mutex
	client    mqttclient.IMqttClient
	connected bool
	lastErr   string
}

// NewPublisher 创建发布器（不建连，由 Start 建立）
func NewPublisher(cfg mqttclient.ClientConfig) *Publisher {
	return &Publisher{cfg: cfg}
}

// Start 建立 MQTT 连接
func (p *Publisher) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client != nil && p.client.IsConnected() {
		return nil
	}
	client := mqttclient.NewMqttClient(p.cfg)
	if client == nil {
		return fmt.Errorf("创建 MQTT 客户端失败，请检查 Broker 配置")
	}
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("连接 MQTT Broker 失败: %w", err)
	}
	p.client = client
	p.connected = true
	p.lastErr = ""
	return nil
}

// Publish 发布消息
func (p *Publisher) Publish(ctx context.Context, topic string, payload []byte, qos byte, retain bool) error {
	p.mu.Lock()
	client := p.client
	connected := p.connected
	p.mu.Unlock()
	if client == nil || !connected {
		return fmt.Errorf("MQTT 未连接")
	}
	opt := mqttclient.PublishOption{
		Topic:   topic,
		Payload: payload,
		Qos:     qos,
		Retain:  retain,
	}
	if err := client.Publish(ctx, opt); err != nil {
		p.mu.Lock()
		p.lastErr = err.Error()
		p.mu.Unlock()
		return err
	}
	return nil
}

// Close 断开 MQTT 连接
func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client != nil {
		_ = p.client.Disconnect()
		p.client = nil
	}
	p.connected = false
}

// Connected 是否已连接
func (p *Publisher) Connected() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.connected && p.client != nil && p.client.IsConnected()
}

// Status 发布器状态
func (p *Publisher) Status() map[string]any {
	p.mu.Lock()
	defer p.mu.Unlock()
	return map[string]any{
		"connected":  p.connected && p.client != nil && p.client.IsConnected(),
		"last_error": p.lastErr,
		"broker":     p.cfg.Broker,
	}
}
