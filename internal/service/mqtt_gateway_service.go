// MQTTW Modbus 网关模块：Wails 绑定层，转发到 internal/modbus 采集核心。
// 提供设备/点位 CRUD、网关启停、在线读取测试与状态查询。
package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"gofly/internal/dao"
	"gofly/internal/dao/model"
	"gofly/internal/modbus"
	"gofly/internal/mqttclient"
	"gofly/internal/utils/gf"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// MqttGatewayService MQTT 网关模块
//
// 采集任务运行在 MQTTW 进程内（方案 A）：进程退出即采集停止，无后台残留。
// 退出前通过 StopGatewayOnExit() 显式停止，保证串口/连接立即释放、MQTT 遗嘱即时发出。
type MqttGatewayService struct{}

// 包级网关实例（Wails 服务单例；退出钩子 StopGatewayOnExit 需要全局可达）
var (
	gwMu       sync.Mutex
	gwInstance *modbus.Gateway
)

// StopGatewayOnExit 进程退出钩子：显式停止采集（幂等，未运行时无操作）。
// 在 main.go 通过 app.OnShutdown 注册。
func StopGatewayOnExit() {
	gwMu.Lock()
	g := gwInstance
	gwMu.Unlock()
	if g != nil && g.Running() {
		g.Stop()
	}
}

// ==================== 设备 CRUD ====================

// ListModbusDevices 设备列表
func (s *MqttGatewayService) ListModbusDevices() any {
	deviceDB := dao.Query().ModbusDevice
	list, err := deviceDB.Order(deviceDB.ID).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取设备列表成功").SetData(list)
}

// SaveModbusDevice 新增/编辑设备（id>0 更新，否则新增）。
// param: 前端 JSON（键名 camelCase 或 snake_case 均可）
func (s *MqttGatewayService) SaveModbusDevice(param string) any {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(param), &data); err != nil {
		return gf.Failed().SetMsg("参数解析失败: " + err.Error())
	}
	id := gf.GetEditId(data["id"])
	norm := normalizeData(data)
	device := &model.ModbusDevice{}
	if b, err := json.Marshal(norm); err == nil {
		_ = json.Unmarshal(b, device)
	}
	deviceDB := dao.Query().ModbusDevice
	if id > 0 {
		delete(norm, "id")
		if _, err := deviceDB.Where(deviceDB.ID.Eq(id)).Updates(norm); err != nil {
			return gf.Failed().SetMsg("更新设备失败: " + err.Error())
		}
		return gf.Success().SetMsg("设备已更新").SetData(map[string]any{"id": id})
	}
	device.ID = 0
	if err := deviceDB.Create(device); err != nil {
		return gf.Failed().SetMsg("新增设备失败: " + err.Error())
	}
	return gf.Success().SetMsg("设备已新增").SetData(map[string]any{"id": device.ID})
}

// DeleteModbusDevice 删除设备（级联删除其点位）
func (s *MqttGatewayService) DeleteModbusDevice(id int32) any {
	if id <= 0 {
		return gf.Failed().SetMsg("无效的设备 ID")
	}
	deviceDB := dao.Query().ModbusDevice
	if _, err := deviceDB.Where(deviceDB.ID.Eq(id)).Delete(); err != nil {
		return gf.Failed().SetMsg("删除设备失败: " + err.Error())
	}
	pointDB := dao.Query().ModbusPoint
	if _, err := pointDB.Where(pointDB.DeviceID.Eq(id)).Delete(); err != nil {
		return gf.Failed().SetMsg("删除设备点位失败: " + err.Error())
	}
	return gf.Success().SetMsg("设备及点位已删除")
}

// ==================== 点位 CRUD ====================

// ListModbusPoints 点位列表
func (s *MqttGatewayService) ListModbusPoints(deviceId int32) any {
	pointDB := dao.Query().ModbusPoint
	list, err := pointDB.Where(pointDB.DeviceID.Eq(deviceId)).Order(pointDB.Sort).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取点位列表成功").SetData(list)
}

// SaveModbusPoint 新增/编辑点位（id>0 更新，否则新增）
func (s *MqttGatewayService) SaveModbusPoint(param string) any {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(param), &data); err != nil {
		return gf.Failed().SetMsg("参数解析失败: " + err.Error())
	}
	id := gf.GetEditId(data["id"])
	norm := normalizeData(data)
	point := &model.ModbusPoint{}
	if b, err := json.Marshal(norm); err == nil {
		_ = json.Unmarshal(b, point)
	}
	if point.DeviceID <= 0 {
		if did, ok := data["deviceId"]; ok {
			point.DeviceID = gf.GetEditId(did)
		}
	}
	pointDB := dao.Query().ModbusPoint
	if id > 0 {
		delete(norm, "id")
		if _, err := pointDB.Where(pointDB.ID.Eq(id)).Updates(norm); err != nil {
			return gf.Failed().SetMsg("更新点位失败: " + err.Error())
		}
		return gf.Success().SetMsg("点位已更新").SetData(map[string]any{"id": id})
	}
	point.ID = 0
	if err := pointDB.Create(point); err != nil {
		return gf.Failed().SetMsg("新增点位失败: " + err.Error())
	}
	return gf.Success().SetMsg("点位已新增").SetData(map[string]any{"id": point.ID})
}

// DeleteModbusPoint 删除点位
func (s *MqttGatewayService) DeleteModbusPoint(id int32) any {
	if id <= 0 {
		return gf.Failed().SetMsg("无效的点位 ID")
	}
	pointDB := dao.Query().ModbusPoint
	if _, err := pointDB.Where(pointDB.ID.Eq(id)).Delete(); err != nil {
		return gf.Failed().SetMsg("删除点位失败: " + err.Error())
	}
	return gf.Success().SetMsg("点位已删除")
}

// BatchImportModbusPoints CSV 批量导入点位。
// CSV 表头（可选，第一行含 name 则跳过）：
//
//	name,register_type,address,quantity,data_type,scale,offset,unit,topic_override,only_on_change
func (s *MqttGatewayService) BatchImportModbusPoints(deviceId int32, csvText string) any {
	if deviceId <= 0 {
		return gf.Failed().SetMsg("请先选择设备")
	}
	reader := csv.NewReader(strings.NewReader(csvText))
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return gf.Failed().SetMsg("CSV 解析失败: " + err.Error())
	}
	if len(rows) == 0 {
		return gf.Failed().SetMsg("CSV 内容为空")
	}
	start := 0
	if strings.Contains(strings.ToLower(rows[0][0]), "name") {
		start = 1
	}
	points := make([]*model.ModbusPoint, 0, len(rows)-start)
	imported := 0
	for i := start; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}
		p := &model.ModbusPoint{
			DeviceID:     deviceId,
			Name:         strings.TrimSpace(row[0]),
			RegisterType: "holding",
			DataType:     "uint16",
			Scale:        1,
			OnlyOnChange: 0,
			Enabled:      1,
		}
		col := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}
		if v := col(1); v != "" {
			p.RegisterType = v
		}
		p.Address = parseInt32(col(2))
		if v := parseInt32(col(3)); v > 0 {
			p.Quantity = v
		}
		if v := col(4); v != "" {
			p.DataType = v
		}
		if v := col(5); v != "" {
			p.Scale = parseFloat(v)
		}
		if v := col(6); v != "" {
			p.Offset_ = parseFloat(v)
		}
		p.Unit = col(7)
		p.TopicOverride = col(8)
		if v := col(9); v == "1" || strings.EqualFold(v, "true") {
			p.OnlyOnChange = 1
		}
		if p.Address < 0 {
			p.Address = 0
		}
		if p.Quantity <= 0 {
			p.Quantity = 1
		}
		points = append(points, p)
		imported++
	}
	if len(points) == 0 {
		return gf.Failed().SetMsg("没有可导入的点位数据")
	}
	pointDB := dao.Query().ModbusPoint
	if err := pointDB.CreateInBatches(points, 200); err != nil {
		return gf.Failed().SetMsg("导入失败: " + err.Error())
	}
	return gf.Success().SetMsg(fmt.Sprintf("成功导入 %d 个点位", imported)).SetData(map[string]any{"imported": imported})
}

// ==================== 网关控制 ====================

// StartModbusGateway 启动网关：所有 enabled 设备进入采集
func (s *MqttGatewayService) StartModbusGateway() any {
	gwMu.Lock()
	defer gwMu.Unlock()
	if gwInstance != nil && gwInstance.Running() {
		return gf.Failed().SetMsg("网关已在运行，请先停止")
	}

	deviceDB := dao.Query().ModbusDevice
	devices, err := deviceDB.Where(deviceDB.Enabled.Eq(1)).Order(deviceDB.ID).Find()
	if err != nil {
		return gf.Failed().SetMsg("读取设备失败: " + err.Error())
	}
	if len(devices) == 0 {
		return gf.Failed().SetMsg("没有 enabled 的设备，请先在设备管理中启用")
	}

	pointDB := dao.Query().ModbusPoint
	points, err := pointDB.Where(pointDB.Enabled.Eq(1)).Find()
	if err != nil {
		return gf.Failed().SetMsg("读取点位失败: " + err.Error())
	}

	cfgs := make([]modbus.DeviceConfig, 0, len(devices))
	pointsByDevice := make(map[int32][]modbus.PointConfig)
	mqttByBroker := make(map[int32]mqttclient.ClientConfig)
	for _, d := range devices {
		cfgs = append(cfgs, deviceConfigOf(d))
	}
	for _, p := range points {
		pointsByDevice[p.DeviceID] = append(pointsByDevice[p.DeviceID], pointConfigOf(p))
	}
	// 按 broker_connection_id 汇总 MQTT 输出配置
	for _, d := range devices {
		if d.BrokerConnectionID <= 0 {
			continue
		}
		if _, ok := mqttByBroker[d.BrokerConnectionID]; ok {
			continue
		}
		conn, err := dao.Query().MqttConnection.Where(dao.Query().MqttConnection.ID.Eq(d.BrokerConnectionID)).First()
		if err != nil {
			continue
		}
		mqttByBroker[d.BrokerConnectionID] = gatewayMQTTConfig(conn)
	}

	g := modbus.NewGateway(func(event, payload string) {
		application.Get().Event.Emit(event, payload)
	})
	if err := g.Start(cfgs, pointsByDevice, mqttByBroker); err != nil {
		return gf.Failed().SetMsg("网关启动失败: " + err.Error())
	}
	gwInstance = g
	return gf.Success().SetMsg("网关已启动").SetData(map[string]any{
		"device_count": len(cfgs),
		"broker_count": len(mqttByBroker),
	})
}

// StopModbusGateway 停止网关
func (s *MqttGatewayService) StopModbusGateway() any {
	gwMu.Lock()
	defer gwMu.Unlock()
	if gwInstance == nil || !gwInstance.Running() {
		return gf.Failed().SetMsg("网关未在运行")
	}
	gwInstance.Stop()
	gwInstance = nil
	return gf.Success().SetMsg("网关已停止")
}

// GetModbusGatewayStatus 网关运行状态快照
func (s *MqttGatewayService) GetModbusGatewayStatus() any {
	gwMu.Lock()
	defer gwMu.Unlock()
	if gwInstance == nil {
		return gf.Success().SetMsg("网关未启动").SetData(map[string]any{
			"running": false, "device_count": 0, "devices": []any{},
		})
	}
	return gf.Success().SetMsg("获取网关状态").SetData(gwInstance.Status())
}

// TestReadModbus 在线读取测试（单次，不启动采集任务）
func (s *MqttGatewayService) TestReadModbus(deviceId int32) any {
	if deviceId <= 0 {
		return gf.Failed().SetMsg("无效的设备 ID")
	}
	deviceDB := dao.Query().ModbusDevice
	device, err := deviceDB.Where(deviceDB.ID.Eq(deviceId)).First()
	if err != nil {
		return gf.Failed().SetMsg("设备不存在: " + err.Error())
	}
	pointDB := dao.Query().ModbusPoint
	points, err := pointDB.Where(pointDB.DeviceID.Eq(deviceId)).Order(pointDB.Sort).Find()
	if err != nil {
		return gf.Failed().SetMsg("读取点位失败: " + err.Error())
	}
	if len(points) == 0 {
		return gf.Failed().SetMsg("该设备没有点位配置")
	}
	gwMu.Lock()
	g := gwInstance
	gwMu.Unlock()
	if g == nil {
		g = modbus.NewGateway(nil)
	}
	result, err := g.TestRead(deviceConfigOf(device), pointConfigsOf(points))
	if err != nil {
		return gf.Failed().SetMsg("测试读取失败: " + err.Error())
	}
	return gf.Success().SetMsg("测试读取完成").SetData(result)
}

// ListGatewayBrokerOptions 可选 MQTT Broker 连接列表（设备表单下拉）
func (s *MqttGatewayService) ListGatewayBrokerOptions() any {
	connectionDB := dao.Query().MqttConnection
	list, err := connectionDB.Where(connectionDB.IsGroup.Eq(0)).Order(connectionDB.ID).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	out := make([]map[string]any, 0, len(list))
	for _, c := range list {
		out = append(out, map[string]any{
			"id":          c.ID,
			"title":       c.Title,
			"host":        c.Host,
			"port":        c.Port,
			"mqttVersion": c.MqttVersion,
			"connected":   c.Connected,
		})
	}
	return gf.Success().SetMsg("获取 Broker 列表成功").SetData(out)
}

// ==================== 转换辅助 ====================

// deviceConfigOf model → 采集配置
func deviceConfigOf(d *model.ModbusDevice) modbus.DeviceConfig {
	return modbus.DeviceConfig{
		ID:                 d.ID,
		Name:               d.Name,
		Protocol:           d.Protocol,
		Host:               d.Host,
		Port:               d.Port,
		SlaveID:            d.SlaveID,
		SerialPort:         d.SerialPort,
		BaudRate:           d.BaudRate,
		DataBits:           d.DataBits,
		StopBits:           d.StopBits,
		Parity:             d.Parity,
		PollInterval:       d.PollInterval,
		Timeout:            d.Timeout,
		RetryCount:         d.RetryCount,
		ByteOrder:          d.ByteOrder,
		WordOrder:          d.WordOrder,
		BrokerConnectionID: d.BrokerConnectionID,
		TopicPrefix:        d.TopicPrefix,
		Qos:                d.Qos,
		Retain:             d.Retain,
		AlarmTopic:         d.AlarmTopic,
	}
}

// pointConfigOf model → 点位配置
func pointConfigOf(p *model.ModbusPoint) modbus.PointConfig {
	return modbus.PointConfig{
		ID:            p.ID,
		DeviceID:      p.DeviceID,
		Name:          p.Name,
		RegisterType:  p.RegisterType,
		Address:       p.Address,
		Quantity:      p.Quantity,
		DataType:      p.DataType,
		Scale:         p.Scale,
		Offset:        p.Offset_,
		Unit:          p.Unit,
		TopicOverride: p.TopicOverride,
		OnlyOnChange:  p.OnlyOnChange,
		Enabled:       p.Enabled,
		Sort:          p.Sort,
	}
}

func pointConfigsOf(ps []*model.ModbusPoint) []modbus.PointConfig {
	out := make([]modbus.PointConfig, 0, len(ps))
	for _, p := range ps {
		out = append(out, pointConfigOf(p))
	}
	return out
}

// gatewayMQTTConfig 由 mqtt_connection 构造网关独立 MQTT 客户端配置
// （ClientID 追加 -gw 后缀，避免与调试连接池客户端冲突）
func gatewayMQTTConfig(dto *model.MqttConnection) mqttclient.ClientConfig {
	props := parseConnectProperties(dto.Properties)
	will := parseWillProperties(dto.Will)
	userProps := parseUserProperties(dto.UserProperties)
	return mqttclient.ClientConfig{
		Broker:                     buildBrokerURL(dto.Protocol, dto.Host, dto.Port, dto.Path, dto.Ssl == 1),
		ClientID:                   dto.ClientID + "-gw",
		Username:                   dto.Username,
		Password:                   dto.Password,
		Protocol:                   mqttclient.ProtocolVersion(dto.MqttVersion),
		KeepAlive:                  uint16(dto.Keepalive),
		WillTopic:                  will.WillTopic,
		WillPayload:                []byte(will.WillPayload),
		WillQos:                    byte(will.WillQos),
		WillRetain:                 will.WillRetain == 1,
		TLSEnable:                  dto.Ssl == 1,
		CleanSession:               dto.Clean == 1,
		CleanStart:                 dto.Clean == 1,
		SessionExpiry:              uint32(props.SessionExpiryInterval),
		AutoReconnect:              dto.Reconnect == 1,
		ReconnectPeriod:            uint16(dto.ReconnectPeriod / 1000),
		ReceiveMaximum:             uint16(props.ReceiveMaximum),
		MaximumPacketSize:          uint32(props.MaximumPacketSize),
		TopicAliasMaximum:          uint16(props.TopicAliasMaximum),
		RequestResponseInformation: props.RequestResponseInformation == 1,
		RequestProblemInformation:  props.RequestProblemInformation == 1,
		UserProperties:             userProps,
	}
}

func parseInt32(s string) int32 {
	v, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	return int32(v)
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}
