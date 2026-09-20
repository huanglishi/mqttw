// 缓冲 channel 任务队列（常驻worker + 有限突发应急worker）
// 目标：消息入口不阻塞，不等待任务完成；限制最大并发；队列满直接丢弃任务，防止雪崩；全部传拷贝数据，不传对象指针
// 并发上限约束（生产重点
// asyncWorkerCount + asyncWorkerBurstMax 总和，要小于等于 GORM MaxOpenConns；
// 如果开启 webhook，总和也要小于 webhook 下游安全并发。
package logicmqtt

import (
	"gofly/internal/dao"
	"gofly/internal/dao/model"
	"log"
	"sync/atomic"
	"time"
)

// 异步任务工作池配置
const (
	asyncWorkerCount     = 16              // 常驻基础后台IO任务(db/webhook)
	asyncWorkerBurstMax  = 16              // 最多额外突发应急worker；总最大并发 = asyncWorkerCount + asyncWorkerBurstMax
	asyncQueueCap        = 4096            // 任务队列最大缓冲
	burstTriggerQueueLen = 1024            // 队列积压超过该阈值，才拉起突发worker
	workerIdleTimeout    = 3 * time.Second // 突发worker空闲多久自动销毁
)

// AsyncJob 异步任务结构体：只存放拷贝后的普通值，禁止放 *mqtt.Client、*packets.Packet指针
type MqttMessageData struct {
	ConnectionID            int32  `json:"connection_id"`
	ClientID                string `json:"client_id"`
	Topic                   string `json:"topic"`
	Payload                 string `json:"payload"`
	Qos                     byte   `json:"qos"`
	Retain                  bool   `json:"retain"`
	SubscriptionIdentifiers string `json:"subscription_identifiers"`
	UserProperties          string `json:"user_properties"`
	ResponseTopic           string `json:"response_topic"`
	ContentType             string `json:"content_type"`
	CorrelationData         string `json:"correlation_data"`
	MessageExpiry           uint32 `json:"message_expiry"`
}

var (
	asyncJobChan       = make(chan *MqttMessageData, asyncQueueCap)
	burstWorkerRunning int32 // 当前运行突发worker数量，atomic安全操作
	dropJobTotal       int64 // 累计丢弃任务计数
)

// bool转数据
func BoolToInt(val bool) int32 {
	if val {
		return 1
	}
	return 0
}
func init() {
	// 程序启动初始化常驻固定worker
	for i := 0; i < asyncWorkerCount; i++ {
		go workerLoop()
	}
	log.Printf("[async-pool] init baseWorker=%d burstMax=%d queueCap=%d", asyncWorkerCount, asyncWorkerBurstMax, asyncQueueCap)
	// defer close(asyncJobChan)
}

// workerLoop 常驻worker，不会主动退出，channel关闭才退出
func workerLoop() {
	for job := range asyncJobChan {
		handleAsyncJob(job)
	}
}

// burstWorkerLoop 突发应急worker：空闲超时自动退出
func burstWorkerLoop() {
	defer atomic.AddInt32(&burstWorkerRunning, -1)

	ticker := time.NewTicker(workerIdleTimeout)
	defer ticker.Stop()

	for {
		select {
		case job, ok := <-asyncJobChan:
			if !ok {
				return
			}
			handleAsyncJob(job)
			ticker.Reset(workerIdleTimeout)

		case <-ticker.C:
			log.Printf("[async] burst worker idle exit")
			return
		}
	}
}

// PushAsyncJob 投递异步任务：永不阻塞调用方，队列满直接丢弃事件；积压达标触发突发worker
func PushAsyncJob(job *MqttMessageData) {
	if job == nil {
		return
	}

	select {
	case asyncJobChan <- job:
		// 队列积压达到阈值，尝试拉起突发worker（硬上限保护）
		if len(asyncJobChan) >= burstTriggerQueueLen {
			running := atomic.LoadInt32(&burstWorkerRunning)
			if running < int32(asyncWorkerBurstMax) {
				if atomic.CompareAndSwapInt32(&burstWorkerRunning, running, running+1) {
					go burstWorkerLoop()
				}
			}
		}

	default:
		atomic.AddInt64(&dropJobTotal, 1)
		log.Printf("[async] queue overflow, drop event dropTotal=%d", atomic.LoadInt64(&dropJobTotal))
	}
}

// 处理异步任务
func handleAsyncJob(job *MqttMessageData) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[hooks async] job panic recovered client=%s err=%v", job.ClientID, r)
		}
	}()
	//数据入库操作
	msgDB := dao.Query().MqttMessage
	msg := model.MqttMessage{
		ConnectionID:            job.ConnectionID,
		MsgType:                 "received",
		Topic:                   job.Topic,
		Payload:                 job.Payload,
		Qos:                     int32(job.Qos),
		Retain:                  BoolToInt(job.Retain),
		SubscriptionIdentifiers: job.SubscriptionIdentifiers,
		UserProperties:          job.UserProperties,
		ResponseTopic:           job.ResponseTopic,
		ContentType:             job.ContentType,
		CorrelationData:         job.CorrelationData,
		MessageExpiry:           int32(job.MessageExpiry),
	}
	// msgDB.Create(&msg)
	// 写库加重试：SQLite 单写者，偶发 SQLITE_BUSY 时等待重试，避免消息丢失
	for i := 0; i < 3; i++ {
		if err := msgDB.Create(&msg); err == nil {
			break
		} else if i < 2 {
			time.Sleep(100 * time.Millisecond)
		} else {
			log.Printf("[async] msg save failed: %v", err.Error())
		}
	}

}

// GetAsyncStats 获取运行时状态，用于定时打印日志/监控埋点
func GetAsyncStats() (queueLen int, burstRunning int32, dropTotal int64) {
	return len(asyncJobChan),
		atomic.LoadInt32(&burstWorkerRunning),
		atomic.LoadInt64(&dropJobTotal)
}

// CloseAsyncPool 优雅关闭，服务退出时调用，关闭通道，worker逐步退出
func CloseAsyncPool() {
	close(asyncJobChan)
	log.Println("[async-pool] closed")
}
