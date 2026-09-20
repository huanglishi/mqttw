// 发布器池：N 个 goroutine 按速率/总量发布，QoS>0 时等待 ack 统计成败
package stress

import (
	"context"
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gonzalop/mq"
)

// Publisher 发布器
type Publisher struct {
	cfg     StressConfig
	metrics *Metrics
	pool    *ClientPool
	payload []byte // 预生成 payload：前8字节 = 发送时间戳(ns) + 填充

	stopCh chan struct{}
	done   chan struct{}
	start  sync.Once
	stop   sync.Once
	// 剩余可发条数（PublishTotal 控制）；-1 表示不限
	remaining atomic.Int64
}

// NewPublisher 创建发布器；payloadSize<8 时仍保留 8 字节时间戳
func NewPublisher(cfg StressConfig, m *Metrics, pool *ClientPool) *Publisher {
	size := cfg.PublishPayloadSize
	if size < 8 {
		size = 8
	}
	payload := make([]byte, size)
	// 填充可打印字符，便于观察
	for i := 8; i < size; i++ {
		payload[i] = byte('a' + i%26)
	}
	return &Publisher{
		cfg:     cfg,
		metrics: m,
		pool:    pool,
		payload: payload,
		stopCh:  make(chan struct{}),
		done:    make(chan struct{}),
	}
}

func (p *Publisher) initRemaining() {
	if p.cfg.PublishTotal > 0 {
		p.remaining.Store(p.cfg.PublishTotal)
	} else {
		p.remaining.Store(-1)
	}
}

// Start 启动发布器池（非阻塞）
func (p *Publisher) Start(ctx context.Context) {
	p.start.Do(func() {
		p.initRemaining()
		var wg sync.WaitGroup
		workers := p.cfg.PublishWorkers
		if workers <= 0 {
			workers = 4
		}
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()
				p.runWorker(ctx)
			}()
		}
		go func() {
			wg.Wait()
			close(p.done)
		}()
	})
}

// runWorker 单个发布循环（按速率令牌桶限速）
func (p *Publisher) runWorker(ctx context.Context) {
	rate := p.cfg.PublishRate
	var ticker *time.Ticker
	if rate > 0 {
		ticker = time.NewTicker(time.Second / time.Duration(rate))
		defer ticker.Stop()
	}
	qos := mq.QoS(p.cfg.PublishQos)
	topic := p.cfg.PublishTopic

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopCh:
			return
		default:
		}
		// 总条数限制
		if rem := p.remaining.Load(); rem == 0 {
			return
		} else if rem > 0 && !p.remaining.CompareAndSwap(rem, rem-1) {
			continue
		}
		// 速率限制
		if ticker != nil {
			<-ticker.C
		}
		client := p.pool.PickClient()
		if client == nil {
			// 无可用连接：短暂等待后重试（等连接阶段完成）
			select {
			case <-ctx.Done():
				return
			case <-time.After(50 * time.Millisecond):
				continue
			}
		}
		p.publishOne(ctx, client, topic, qos)
	}
}

// publishOne 单条发布；写时间戳进 payload 前 8 字节
func (p *Publisher) publishOne(ctx context.Context, client *mq.Client, topic string, qos mq.QoS) {
	binary.LittleEndian.PutUint64(p.payload[:8], uint64(time.Now().UnixNano()))
	token := client.Publish(ctx, topic, p.payload, mq.WithQoS(qos))
	p.metrics.PublishTotal.Add(1)
	if qos == mq.AtMostOnce {
		p.metrics.PublishSuccess.Add(1) // QoS0 即发即成功
		return
	}
	if err := token.Wait(ctx); err != nil {
		p.metrics.PublishFail.Add(1)
		p.metrics.AddFailReason("publish:" + err.Error())
		return
	}
	p.metrics.PublishSuccess.Add(1)
}

// Stop 停止发布（等待 worker 退出）
func (p *Publisher) Stop() {
	p.stop.Do(func() {
		close(p.stopCh)
		select {
		case <-p.done:
		case <-time.After(3 * time.Second):
		}
	})
}
