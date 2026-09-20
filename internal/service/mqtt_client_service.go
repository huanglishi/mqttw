package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gofly/internal/dao"
	logicmqtt "gofly/internal/logic/logic_mqtt"
	"gofly/internal/mqttclient"
	"gofly/internal/utils/gf"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// 客户端连接MQTT-Broker操作
type MqttClientService struct {
	clientMap map[string]mqttclient.IMqttClient
	mu        sync.RWMutex
}

// func NewMqttService(app *application.App) *MqttClientService {
// 	return &MqttClientService{
// 		app:       app,
// 		clientMap: make(map[string]mqttclient.IMqttClient),
// 	}
// }

// ==================== Connect 内部辅助（仅本文件使用） ====================

// ConnectProperties 对应 mqtt_connection.properties 列存储的 JSON（MQTT5 CONNECT 属性）
// 前端 a-switch 的 checked-value=1 / unchecked-value=0，故布尔类字段用 int32 接收
type ConnectProperties struct {
	SessionExpiryInterval      int32 `json:"sessionExpiryInterval"`
	ReceiveMaximum             int32 `json:"receiveMaximum"`
	MaximumPacketSize          int32 `json:"maximumPacketSize"`
	TopicAliasMaximum          int32 `json:"topicAliasMaximum"`
	RequestResponseInformation int32 `json:"requestResponseInformation"`
	RequestProblemInformation  int32 `json:"requestProblemInformation"`
}

// WillProperties 对应 mqtt_connection.will 列存储的 JSON（遗嘱消息）
type WillProperties struct {
	WillTopic              string `json:"willTopic"`
	WillQos                int32  `json:"willQos"`
	WillRetain             int32  `json:"willRetain"`
	WillPayload            string `json:"willPayload"`
	PayloadFormatIndicator int32  `json:"payloadFormatIndicator"`
	WillDelayInterval      int32  `json:"willDelayInterval"`
	MessageExpiryInterval  int32  `json:"messageExpiryInterval"`
	ContentType            string `json:"contentType"`
	ResponseTopic          string `json:"responseTopic"`
	CorrelationData        string `json:"correlationData"`
}

// UserProperty 对应 mqtt_connection.user_properties 列存储的 JSON 数组元素（KeyValueEditor 的 v-model 格式）
type UserProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// parseConnectProperties 解析 properties JSON，失败返回零值
func parseConnectProperties(raw string) ConnectProperties {
	var p ConnectProperties
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &p)
	}
	return p
}

// parseWillProperties 解析 will JSON，失败返回零值
func parseWillProperties(raw string) WillProperties {
	var w WillProperties
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &w)
	}
	return w
}

// parseUserProperties 解析 user_properties JSON [{key,value}] → mqttclient 需要的 [][]string
func parseUserProperties(raw string) [][]string {
	if raw == "" {
		return nil
	}
	var list []UserProperty
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return nil
	}
	var out [][]string
	for _, item := range list {
		if item.Key == "" {
			continue // 空 key 不携带
		}
		out = append(out, []string{item.Key, item.Value})
	}
	return out
}

// buildBrokerURL 由前端协议字段构造 paho 可识别的 Broker URL
// 前端 protocol 取值：mqtt / mqtts / ws / wss（无 :// 后缀）；ssl=1 时 mqtt→ssl、ws→wss
func buildBrokerURL(protocol, host, port, path string, ssl bool) string {
	scheme := strings.ToLower(protocol)
	switch scheme {
	case "mqtt":
		scheme = "tcp"
	case "mqtts":
		scheme = "ssl"
	case "ws", "wss":
	default:
		scheme = "tcp"
	}

	if ssl {
		switch scheme {
		case "tcp":
			scheme = "ssl"
		case "ws":
			scheme = "wss"
		}
	}
	broker := fmt.Sprintf("%s://%s:%s", scheme, host, port)
	if (scheme == "ws" || scheme == "wss") && path != "" {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		broker += path
	}
	return broker
}

// Connect 客户端连接MQTT服务器
func (s *MqttClientService) Connect(id int32) any {
	// 1. 获取连接数据
	connectionDB := dao.Query().MqttConnection
	dto, err := connectionDB.Where(connectionDB.ID.Eq(id)).First()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}

	// 2. 解析 DB 中 TEXT 列存储的 JSON（properties / will / user_properties）
	props := parseConnectProperties(dto.Properties)
	will := parseWillProperties(dto.Will)
	userProps := parseUserProperties(dto.UserProperties)

	// 3. 构造客户端配置
	cfg := mqttclient.ClientConfig{
		Broker:                     buildBrokerURL(dto.Protocol, dto.Host, dto.Port, dto.Path, dto.Ssl == 1),
		ClientID:                   dto.ClientID,
		Username:                   dto.Username,
		Password:                   dto.Password,
		Protocol:                   mqttclient.ProtocolVersion(dto.MqttVersion),
		KeepAlive:                  uint16(dto.Keepalive),
		WillTopic:                  will.WillTopic,
		WillPayload:                []byte(will.WillPayload),
		WillQos:                    byte(will.WillQos),
		WillRetain:                 will.WillRetain == 1,
		TLSEnable:                  dto.Ssl == 1,
		CleanSession:               dto.Clean == 1, // MQTT3.x
		CleanStart:                 dto.Clean == 1, // MQTT5
		SessionExpiry:              uint32(props.SessionExpiryInterval),
		AutoReconnect:              dto.Reconnect == 1,
		ReconnectPeriod:            uint16(dto.ReconnectPeriod / 1000), // 前端毫秒 → 底层秒
		ReceiveMaximum:             uint16(props.ReceiveMaximum),
		MaximumPacketSize:          uint32(props.MaximumPacketSize),
		TopicAliasMaximum:          uint16(props.TopicAliasMaximum),
		RequestResponseInformation: props.RequestResponseInformation == 1,
		RequestProblemInformation:  props.RequestProblemInformation == 1,
		UserProperties:             userProps,
	}

	// 4. 创建客户端（按 MQTT 版本分派：3.1/3.1.1 老 Paho，5.0 autopaho）
	client := mqttclient.NewMqttClient(cfg)
	if client == nil {
		return gf.Failed().SetMsg("创建MQTT客户端失败，请检查 MQTT 版本或 Broker 地址")
	}

	// 5. 消息回调，事件推送到前端
	client.OnMessage(func(msg *mqttclient.MqttMessage) {
		// MqttMessage 转 JSON 字符串（Payload/CorrelationData 转可读文本，避免 base64）
		out := &logicmqtt.MqttMessageData{
			ConnectionID:            dto.ID,
			ClientID:                dto.ClientID,
			Topic:                   msg.Topic,
			Payload:                 string(msg.Payload),
			Qos:                     msg.Qos,
			Retain:                  msg.Retain,
			SubscriptionIdentifiers: gf.Uint32SliceToString(msg.SubscriptionIdentifiers),
			UserProperties:          gf.DoubleSliceToString(msg.UserProperties),
			ResponseTopic:           msg.ResponseTopic,
			ContentType:             msg.ContentType,
			CorrelationData:         string(msg.CorrelationData),
			MessageExpiry:           msg.MessageExpiry,
		}
		logicmqtt.PushAsyncJob(out)
		data, err := json.Marshal(out)
		if err != nil {
			return
		}
		app := application.Get()
		app.Event.Emit("mqtt:all_message", string(data))
		app.Event.Emit("mqtt:message:26", string(data))
	})

	// 6. 连接前先断开池中同 ID 旧客户端：
	key := strconv.Itoa(int(dto.ID))
	s.mu.Lock()
	if s.clientMap == nil {
		s.clientMap = make(map[string]mqttclient.IMqttClient)
	}
	if old, ok := s.clientMap[key]; ok {
		_ = old.Disconnect()
		delete(s.clientMap, key)
	}
	s.mu.Unlock()

	// 7. 连接，超时取连接配置（默认 10s）
	timeout := time.Duration(dto.ConnectTimeout) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return gf.Failed().SetMsg("连接失败：" + err.Error())
	}

	//8.在连接成功后把连接下的订阅主题数据列表订阅上
	subscriptionDB := dao.Query().MqttSubscription
	subscription, err := subscriptionDB.Where(subscriptionDB.Disabled.Eq(0), subscriptionDB.ConnectionID.Eq(id)).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	var failedTopics []string
	if len(subscription) > 0 {
		for _, sub := range subscription {
			if sub.Topic == "" {
				continue // 跳过空主题
			}
			opt := mqttclient.SubscribeOption{
				Topic:                  sub.Topic,
				Qos:                    byte(sub.Qos),
				SubscriptionIdentifier: uint32(sub.SubscriptionIdentifier),
				NoLocal:                sub.Nl == 1,
				RetainAsPublished:      sub.Rap == 1,
				RetainHandling:         byte(sub.Rh),
			}
			subCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := client.Subscribe(subCtx, opt); err != nil {
				failedTopics = append(failedTopics, sub.Topic)
			}
			cancel()
		}
	}

	// 9. 存入连接池
	s.mu.Lock()
	if s.clientMap == nil {
		s.clientMap = make(map[string]mqttclient.IMqttClient)
	}
	s.clientMap[key] = client
	s.mu.Unlock()
	//10.更新连接数据连接状态
	connectionDB.Where(connectionDB.ID.Eq(id)).Update(connectionDB.Connected, 1)
	return gf.Success().SetMsg("客户端连接成功").SetData(cfg)
}

// Disconnect 客户端断开连接
func (s *MqttClientService) Disconnect(id int32) any {
	key := strconv.Itoa(int(id))
	s.mu.Lock()
	if old, ok := s.clientMap[key]; ok {
		if err := old.Disconnect(); err != nil {
			s.mu.Unlock()
			return gf.Failed().SetMsg("断开连接失败：" + err.Error())
		}
		delete(s.clientMap, key)
	}
	s.mu.Unlock()
	// 更新连接数据连接状态
	connectionDB := dao.Query().MqttConnection
	connectionDB.Where(connectionDB.ID.Eq(id)).Update(connectionDB.Connected, 0)
	return gf.Success().SetMsg("客户端断开连接成功")
}

// DisconnectAll 断开全部在线客户端（全局快捷键 Ctrl+Shift+D 调用）
func (s *MqttClientService) DisconnectAll() any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.clientMap) == 0 {
		return gf.Success().SetMsg("当前没有在线连接")
	}
	count := 0
	connectionDB := dao.Query().MqttConnection
	for key, client := range s.clientMap {
		if client == nil {
			continue
		}
		_ = client.Disconnect()
		delete(s.clientMap, key)
		if id, err := strconv.Atoi(key); err == nil {
			connectionDB.Where(connectionDB.ID.Eq(int32(id))).Update(connectionDB.Connected, 0)
		}
		count++
	}
	return gf.Success().SetMsg(fmt.Sprintf("已断开全部连接（%d 个）", count)).SetData(count)
}

// getClient 从连接池取客户端，未连接返回错误
func (s *MqttClientService) getClient(id int32) (mqttclient.IMqttClient, error) {
	key := strconv.Itoa(int(id))
	s.mu.RLock()
	client, ok := s.clientMap[key]
	s.mu.RUnlock()
	if !ok || client == nil {
		return nil, fmt.Errorf("连接不存在或未连接，请先连接")
	}
	return client, nil
}

// rawToBool 兼容前端 bool（true/false）与 0/1 数字两种布尔表达
func rawToBool(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var b bool
	if json.Unmarshal(raw, &b) == nil {
		return b
	}
	var n int
	if json.Unmarshal(raw, &n) == nil {
		return n != 0
	}
	return false
}

// Subscribe 订阅主题
func (s *MqttClientService) Subscribe(id int32) any {
	subscriptionDB := dao.Query().MqttSubscription
	req, err := subscriptionDB.Where(subscriptionDB.ID.Eq((id))).First()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}

	if req.ConnectionID <= 0 || req.Topic == "" {
		return gf.Failed().SetMsg("connection_id 和 topic 不能为空")
	}
	client, err := s.getClient(req.ConnectionID)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	opt := mqttclient.SubscribeOption{
		Topic:                  req.Topic,
		Qos:                    byte(req.Qos),
		SubscriptionIdentifier: uint32(req.SubscriptionIdentifier),
		NoLocal:                req.Nl == 1,
		RetainAsPublished:      req.Rap == 1,
		RetainHandling:         byte(req.Rh),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Subscribe(ctx, opt); err != nil {
		return gf.Failed().SetMsg("订阅失败：" + err.Error())
	}
	return gf.Success().SetMsg("客户端订阅主题成功").SetData(opt)
}

// UnSubscribe 取消主题订阅
// param JSON 示例：{"connection_id":1,"topic":"topic/#"}
func (s *MqttClientService) UnSubscribe(connection_id int32, topic string) any {
	if connection_id <= 0 || topic == "" {
		return gf.Failed().SetMsg("connection_id 和 topic 不能为空")
	}
	client, err := s.getClient(connection_id)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.UnSubscribe(ctx, topic); err != nil {
		return gf.Failed().SetMsg("取消订阅失败：" + err.Error())
	}
	return gf.Success().SetMsg("取消主题订阅成功")
}

// Publish 发布消息
// param JSON 示例（payload 支持字符串或 JSON 对象；retain 支持 bool 或 0/1）：
//
//	{"connection_id":1,"topic":"t","payload":"{\"k\":1}","qos":0,"retain":false,
//	 "user_properties":[{"key":"k","value":"v"}],"response_topic":"","content_type":"",
//	 "correlation_data":"","message_expiry":0}
func (s *MqttClientService) Publish(param string) any {
	var req struct {
		ConnectionID           int32           `json:"connection_id"`
		Topic                  string          `json:"topic"`
		Payload                json.RawMessage `json:"payload"`
		Qos                    int32           `json:"qos"`
		Retain                 json.RawMessage `json:"retain"`
		UserProperties         json.RawMessage `json:"user_properties"`
		ResponseTopic          string          `json:"response_topic"`
		ContentType            string          `json:"content_type"`
		CorrelationData        string          `json:"correlation_data"`
		MessageExpiry          int64           `json:"message_expiry"`
		TopicAlias             int64           `json:"topic_alias"`
		PayloadFormat          bool            `json:"payload_format"`
		SubscriptionIdentifier int             `json:"subscription_identifier"`
	}

	if err := json.Unmarshal([]byte(param), &req); err != nil {
		return gf.Failed().SetMsg("发布参数解析失败，" + err.Error())
	}
	if req.ConnectionID <= 0 || req.Topic == "" {
		return gf.Failed().SetMsg("connection_id 和 topic 不能为空")
	}
	client, err := s.getClient(req.ConnectionID)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}

	// payload：字符串原样使用，对象/数组使用原始 JSON 字节
	var payload []byte
	if len(req.Payload) > 0 {
		var s string
		if json.Unmarshal(req.Payload, &s) == nil {
			payload = []byte(s)
		} else {
			payload = req.Payload
		}
	}

	opt := mqttclient.PublishOption{
		Topic:                  req.Topic,
		Payload:                payload,
		Qos:                    byte(req.Qos),
		Retain:                 rawToBool(req.Retain),
		UserProperties:         parseUserProperties(string(req.UserProperties)),
		ResponseTopic:          req.ResponseTopic,
		ContentType:            req.ContentType,
		CorrelationData:        []byte(req.CorrelationData),
		MessageExpiry:          uint32(req.MessageExpiry),
		TopicAlias:             uint16(req.TopicAlias),
		PayloadFormat:          byte(logicmqtt.BoolToInt(req.PayloadFormat)),
		SubscriptionIdentifier: req.SubscriptionIdentifier,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Publish(ctx, opt); err != nil {
		return gf.Failed().SetMsg("发布失败：" + err.Error())
	}
	return gf.Success().SetMsg("客户端发布消息成功").SetData(opt)
}
