// MQTT5.0 Paho.golang autopaho 客户端实现（适配 paho.golang v0.23.0 API）
package mqttclient

import (
	"context"
	"crypto/tls"
	"net/url"
	"sync"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

const connectTimeout = 10 * time.Second

type PahoV5Client struct {
	cfg        ClientConfig
	cm         *autopaho.ConnectionManager
	msgHandler func(msg *MqttMessage)
	mu         sync.Mutex
	subs       []*paho.Subscribe // 保存订阅，重连自动恢复

	connMu    sync.Mutex
	connected bool
}

func NewPahoV5Client(cfg ClientConfig) *PahoV5Client {
	cli := &PahoV5Client{cfg: cfg}

	// Broker URL 解析，支持 mqtt:// tcp:// ws:// wss:// 等 scheme
	u, err := url.Parse(cfg.Broker)
	if err != nil {
		return nil
	}

	apCfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     cfg.KeepAlive,
		CleanStartOnInitialConnection: cfg.CleanStart,
		SessionExpiryInterval:         cfg.SessionExpiry,
		ConnectUsername:               cfg.Username,
		ConnectPassword:               []byte(cfg.Password),
		ConnectTimeout:                connectTimeout,
		ReconnectBackoff: func(int) time.Duration {
			return time.Duration(cfg.ReconnectPeriod) * time.Second
		},
		OnConnectionUp: func(cm *autopaho.ConnectionManager, _ *paho.Connack) {
			cli.connMu.Lock()
			cli.connected = true
			cli.connMu.Unlock()

			// 重连成功自动恢复所有订阅
			cli.mu.Lock()
			defer cli.mu.Unlock()
			for _, sub := range cli.subs {
				_, _ = cm.Subscribe(context.Background(), sub)
			}
		},
		OnConnectionDown: func() bool {
			cli.connMu.Lock()
			cli.connected = false
			cli.connMu.Unlock()
			// 返回 false 时停止重连；开启自动重连则继续尝试
			return cfg.AutoReconnect
		},
		OnConnectError: func(error) {},
		ClientConfig: paho.ClientConfig{
			ClientID: cfg.ClientID,
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					if cli.msgHandler != nil {
						cli.msgHandler(msgFromPacket(pr.Packet))
					}
					return true, nil
				},
			},
			OnClientError: func(error) {},
		},
	}

	// MQTT5 CONNECT 属性：ReceiveMaximum / MaximumPacketSize / TopicAliasMaximum /
	// RequestResponseInfo / RequestProblemInfo / 连接用户属性
	apCfg.ConnectPacketBuilder = func(c *paho.Connect, _ *url.URL) (*paho.Connect, error) {
		c.Properties = &paho.ConnectProperties{
			ReceiveMaximum:      ptrUint16(cfg.ReceiveMaximum),
			MaximumPacketSize:   ptrUint32(cfg.MaximumPacketSize),
			TopicAliasMaximum:   ptrUint16(cfg.TopicAliasMaximum),
			RequestResponseInfo: cfg.RequestResponseInformation,
			RequestProblemInfo:  cfg.RequestProblemInformation,
			User:                toPahoUserProps(cfg.UserProperties),
		}
		return c, nil
	}

	// 遗嘱消息
	if cfg.WillTopic != "" {
		apCfg.WillMessage = &paho.WillMessage{
			Topic:   cfg.WillTopic,
			Payload: cfg.WillPayload,
			QoS:     cfg.WillQos,
			Retain:  cfg.WillRetain,
		}
	}

	if cfg.TLSEnable {
		apCfg.TlsCfg = &tls.Config{InsecureSkipVerify: true}
	}

	cm, err := autopaho.NewConnection(context.Background(), apCfg)
	if err != nil {
		return nil
	}
	cli.cm = cm
	return cli
}

func (c *PahoV5Client) Connect(ctx context.Context) error {
	return c.cm.AwaitConnection(ctx)
}

func (c *PahoV5Client) Disconnect() error {
	return c.cm.Disconnect(context.Background())
}

func (c *PahoV5Client) Subscribe(ctx context.Context, opt SubscribeOption) error {
	sub := &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{
				Topic:             opt.Topic,
				QoS:               opt.Qos,
				RetainHandling:    opt.RetainHandling,
				NoLocal:           opt.NoLocal,
				RetainAsPublished: opt.RetainAsPublished,
			},
		},
	}

	// MQTT5 订阅标识符（0 时不携带该属性）
	if opt.SubscriptionIdentifier > 0 {
		id := int(opt.SubscriptionIdentifier)
		sub.Properties = &paho.SubscribeProperties{
			SubscriptionIdentifier: &id,
		}
	}

	c.mu.Lock()
	c.subs = append(c.subs, sub)
	c.mu.Unlock()

	_, err := c.cm.Subscribe(ctx, sub)
	return err
}

func (c *PahoV5Client) UnSubscribe(ctx context.Context, topic string) error {
	unsub := &paho.Unsubscribe{Topics: []string{topic}}
	_, err := c.cm.Unsubscribe(ctx, unsub)
	if err != nil {
		return err
	}

	// 删除本地缓存订阅
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, s := range c.subs {
		if len(s.Subscriptions) > 0 && s.Subscriptions[0].Topic == topic {
			c.subs = append(c.subs[:i], c.subs[i+1:]...)
			break
		}
	}
	return nil
}

func (c *PahoV5Client) Publish(ctx context.Context, opt PublishOption) error {
	pub := &paho.Publish{
		Topic:   opt.Topic,
		Payload: opt.Payload,
		QoS:     opt.Qos,
		Retain:  opt.Retain,
	}

	props := &paho.PublishProperties{
		User:            toPahoUserProps(opt.UserProperties),
		ResponseTopic:   opt.ResponseTopic,
		ContentType:     opt.ContentType,
		CorrelationData: opt.CorrelationData,
	}
	// MQTT5 PUBLISH 属性（0 时不携带该属性）
	if opt.TopicAlias > 0 {
		props.TopicAlias = &opt.TopicAlias
	}
	if opt.PayloadFormat > 0 {
		props.PayloadFormat = &opt.PayloadFormat
	}
	if opt.MessageExpiry > 0 {
		props.MessageExpiry = &opt.MessageExpiry
	}
	if opt.SubscriptionIdentifier > 0 {
		props.SubscriptionIdentifier = &opt.SubscriptionIdentifier
	}
	pub.Properties = props

	_, err := c.cm.Publish(ctx, pub)
	return err
}

func (c *PahoV5Client) OnMessage(cb func(msg *MqttMessage)) {
	c.msgHandler = cb
}

func (c *PahoV5Client) IsConnected() bool {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	return c.connected
}

// msgFromPacket 将 paho.Publish 转换为统一 MqttMessage
func msgFromPacket(p *paho.Publish) *MqttMessage {
	msg := &MqttMessage{
		Topic:   p.Topic,
		Payload: p.Payload,
		Qos:     p.QoS,
		Retain:  p.Retain,
	}
	if p.Properties == nil {
		return msg
	}
	if p.Properties.MessageExpiry != nil {
		msg.MessageExpiry = *p.Properties.MessageExpiry
	}
	msg.ResponseTopic = p.Properties.ResponseTopic
	msg.ContentType = p.Properties.ContentType
	msg.CorrelationData = p.Properties.CorrelationData
	if p.Properties.SubscriptionIdentifier != nil {
		msg.SubscriptionIdentifiers = []uint32{uint32(*p.Properties.SubscriptionIdentifier)}
	}
	for _, up := range p.Properties.User {
		msg.UserProperties = append(msg.UserProperties, []string{up.Key, up.Value})
	}
	return msg
}

// toPahoUserProps 将 [][]string 键值对转为 paho.UserProperties
func toPahoUserProps(props [][]string) paho.UserProperties {
	if len(props) == 0 {
		return nil
	}
	var out paho.UserProperties
	for _, kv := range props {
		if len(kv) >= 2 {
			out = append(out, paho.UserProperty{Key: kv[0], Value: kv[1]})
		}
	}
	return out
}

func ptrUint16(v uint16) *uint16 {
	if v == 0 {
		return nil
	}
	return &v
}

func ptrUint32(v uint32) *uint32 {
	if v == 0 {
		return nil
	}
	return &v
}
