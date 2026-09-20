// MQTT3.1 / MQTT3.1.1 老 Paho 客户端实现
package mqttclient

import (
	"context"
	"crypto/tls"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PahoClient struct {
	cfg        ClientConfig
	client     mqtt.Client
	msgHandler func(msg *MqttMessage)
}

func NewPahoClient(cfg ClientConfig) *PahoClient {
	c := &PahoClient{cfg: cfg}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)
	opts.SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second)

	// 区分 3.1 / 3.1.1 协议版本
	switch cfg.Protocol {
	case ProtocolV31:
		opts.SetProtocolVersion(3) // MQIsdp ver3
	case ProtocolV311:
		opts.SetProtocolVersion(4) // MQTT ver4
	}

	opts.SetCleanSession(cfg.CleanSession)
	opts.SetAutoReconnect(cfg.AutoReconnect)
	if cfg.ReconnectPeriod > 0 {
		opts.SetConnectRetryInterval(time.Duration(cfg.ReconnectPeriod) * time.Second)
	}

	// 遗嘱消息（老 Paho SetWill 的 payload 是 string）
	if cfg.WillTopic != "" {
		opts.SetWill(cfg.WillTopic, string(cfg.WillPayload), cfg.WillQos, cfg.WillRetain)
	}

	if cfg.TLSEnable {
		opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	// 默认消息回调，注意 MessageHandler 签名为 func(Client, Message)
	opts.SetDefaultPublishHandler(func(_ mqtt.Client, pahoMsg mqtt.Message) {
		if c.msgHandler == nil {
			return
		}
		msg := &MqttMessage{
			Topic:                   pahoMsg.Topic(),
			Payload:                 pahoMsg.Payload(),
			Qos:                     pahoMsg.Qos(),
			Retain:                  pahoMsg.Retained(),
			SubscriptionIdentifiers: nil,
			UserProperties:          nil,
		}
		c.msgHandler(msg)
	})

	c.client = mqtt.NewClient(opts)
	return c
}

func (p *PahoClient) Connect(ctx context.Context) error {
	token := p.client.Connect()
	if !token.WaitTimeout(connectTimeout) {
		return context.DeadlineExceeded
	}
	return token.Error()
}

func (p *PahoClient) Disconnect() error {
	if p.client.IsConnected() {
		p.client.Disconnect(250)
	}
	return nil
}

func (p *PahoClient) Subscribe(ctx context.Context, opt SubscribeOption) error {
	// 第三个参数传 nil 使用 DefaultPublishHandler 接收消息
	token := p.client.Subscribe(opt.Topic, opt.Qos, nil)
	if !token.WaitTimeout(connectTimeout) {
		return context.DeadlineExceeded
	}
	return token.Error()
}

func (p *PahoClient) UnSubscribe(ctx context.Context, topic string) error {
	token := p.client.Unsubscribe(topic)
	if !token.WaitTimeout(connectTimeout) {
		return context.DeadlineExceeded
	}
	return token.Error()
}

func (p *PahoClient) Publish(ctx context.Context, opt PublishOption) error {
	token := p.client.Publish(opt.Topic, opt.Qos, opt.Retain, opt.Payload)
	if !token.WaitTimeout(connectTimeout) {
		return context.DeadlineExceeded
	}
	return token.Error()
}

func (p *PahoClient) OnMessage(cb func(msg *MqttMessage)) {
	p.msgHandler = cb
}

func (p *PahoClient) IsConnected() bool {
	return p.client.IsConnected()
}
