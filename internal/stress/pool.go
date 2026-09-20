// 客户端池：批量创建/连接/订阅/销毁 gonzalop/mq 客户端
// gonzalop/mq 无 per-connection goroutine，可支撑数千~数万连接
package stress

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gonzalop/mq"
)

// poolClient 单个压测客户端
type poolClient struct {
	id     string // clientID
	client *mq.Client
	subbed int // 实际订阅数
}

// ClientPool 客户端池
type ClientPool struct {
	cfg     StressConfig
	metrics *Metrics
	clients []*poolClient // 成功连接的客户端
	mu      sync.Mutex
	closed  bool
}

// NewClientPool 创建客户端池
func NewClientPool(cfg StressConfig, m *Metrics) *ClientPool {
	return &ClientPool{cfg: cfg, metrics: m, clients: make([]*poolClient, 0, cfg.ClientCount)}
}

// protocolVersion 前端协议 → gonzalop 协议版本
// gonzalop/mq 最低支持 3.1.1；"3.1" 选项以 3.1.1 兼容模式连接（绝大多数 3.1 broker 接受 3.1.1 客户端）
func (p *ClientPool) protocolVersion() (uint8, bool) {
	switch p.cfg.Protocol {
	case "5.0":
		return mq.ProtocolV50, true
	case "3.1", "3.1.1":
		return mq.ProtocolV311, true
	default: // auto：由库自动协商
		return mq.ProtocolV50, false
	}
}

// dialOptions 构建连接选项
func (p *ClientPool) dialOptions(id string) []mq.Option {
	opts := []mq.Option{
		mq.WithClientID(id),
		mq.WithKeepAlive(time.Duration(p.cfg.KeepAlive) * time.Second),
		mq.WithCleanSession(p.cfg.CleanSession),
		mq.WithConnectTimeout(p.cfg.connectTimeoutDuration()),
	}
	if ver, fixed := p.protocolVersion(); fixed {
		opts = append(opts, mq.WithProtocolVersion(ver))
	} else {
		opts = append(opts, mq.WithAutoProtocolVersion(true))
	}
	if p.cfg.Username != "" || p.cfg.Password != "" {
		opts = append(opts, mq.WithCredentials(p.cfg.Username, p.cfg.Password))
	}
	return opts
}

// Start 批量连接 + 订阅（阻塞直到全部完成或 ctx 取消）
// 返回成功连接数；失败原因进入 metrics
func (p *ClientPool) Start(ctx context.Context) int {
	interval := p.cfg.connectInterval()
	var ticker *time.Ticker
	if interval > 0 {
		ticker = time.NewTicker(interval)
		defer ticker.Stop()
	}

	connected := 0
	for i := 0; i < p.cfg.ClientCount; i++ {
		select {
		case <-ctx.Done():
			return connected
		default:
		}
		if ticker != nil {
			<-ticker.C
		}

		id := fmt.Sprintf("%s_%d", p.cfg.ClientPrefix, i+1)
		start := time.Now()
		client, err := mq.DialContext(ctx, p.cfg.Broker, p.dialOptions(id)...)
		if err != nil {
			p.metrics.ConnectTotal.Add(1)
			p.metrics.ConnectFail.Add(1)
			p.metrics.AddFailReason(classifyConnectErr(err))
			continue
		}
		p.metrics.AddConnectTime(time.Since(start))
		p.metrics.ConnectTotal.Add(1)
		p.metrics.ConnectSuccess.Add(1)

		pc := &poolClient{id: id, client: client}
		// 订阅
		if p.cfg.SubscribePerClient > 0 && len(p.cfg.Topics) > 0 {
			p.subscribe(ctx, pc)
		}
		p.mu.Lock()
		if p.closed {
			p.mu.Unlock()
			_ = client.Disconnect(context.Background())
			return connected
		}
		p.clients = append(p.clients, pc)
		p.mu.Unlock()
		connected++
	}
	return connected
}

// subscribe 单个客户端订阅 SubscribePerClient 个 topic（循环取模板）
func (p *ClientPool) subscribe(ctx context.Context, pc *poolClient) {
	qos := mq.QoS(p.cfg.Qos)
	for i := 0; i < p.cfg.SubscribePerClient; i++ {
		topic := p.cfg.Topics[i%len(p.cfg.Topics)]
		token := pc.client.Subscribe(ctx, topic, qos, func(c *mq.Client, msg mq.Message) {
			p.onMessage(msg)
		})
		p.metrics.SubscribeTotal.Add(1)
		if err := token.Wait(ctx); err != nil {
			p.metrics.SubscribeFail.Add(1)
			p.metrics.AddFailReason("subscribe:" + topic + ":" + err.Error())
			continue
		}
		p.metrics.SubscribeSuccess.Add(1)
		pc.subbed++
	}
}

// onMessage 统一消息回调：接收计数 + 回显延迟
func (p *ClientPool) onMessage(msg mq.Message) {
	p.metrics.ReceivedTotal.Add(1)
	if p.cfg.EchoEnabled && msg.Topic == p.cfg.PublishEchoTopic && len(msg.Payload) >= 8 {
		sentAt := int64(msg.Payload[0]) | int64(msg.Payload[1])<<8 | int64(msg.Payload[2])<<16 | int64(msg.Payload[3])<<24 |
			int64(msg.Payload[4])<<32 | int64(msg.Payload[5])<<40 | int64(msg.Payload[6])<<48 | int64(msg.Payload[7])<<56
		p.metrics.AddEchoLatency(time.Since(time.Unix(0, sentAt)))
	}
}

// ConnectedCount 当前已连接客户端数（供快照）
func (p *ClientPool) ConnectedCount() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	var n int64
	for _, pc := range p.clients {
		if pc.client != nil && pc.client.IsConnected() {
			n++
		}
	}
	return n
}

// PickClient 取一个已连接客户端用于发布；无则返回 nil
func (p *ClientPool) PickClient() *mq.Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.clients) == 0 {
		return nil
	}
	for _, pc := range p.clients {
		if pc.client != nil && pc.client.IsConnected() {
			return pc.client
		}
	}
	return nil
}

// Close 断开全部客户端
func (p *ClientPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return
	}
	p.closed = true
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, pc := range p.clients {
		if pc.client != nil {
			_ = pc.client.Disconnect(ctx)
		}
	}
	p.clients = nil
}

// classifyConnectErr 连接错误粗分类（供前端展示）
func classifyConnectErr(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "timeout"):
		return "connect timeout"
	case strings.Contains(lower, "refused"):
		return "connection refused"
	case strings.Contains(lower, "not authorized"), strings.Contains(lower, "bad user"):
		return "auth rejected"
	case strings.Contains(lower, "unacceptable"):
		return "unacceptable protocol/clientID"
	case strings.Contains(lower, "server busy"), strings.Contains(lower, "quota"):
		return "server busy/quota"
	default:
		if len(msg) > 60 {
			msg = msg[:60]
		}
		return msg
	}
}
