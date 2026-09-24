package service

import (
	"encoding/json"

	"gofly/internal/dao"
	"gofly/internal/dao/model"
	"gofly/internal/utils/gf"
)

type MqttMessageRecordService struct{}

// GetList 获取消息记录列表
func (g *MqttMessageRecordService) GetList() any {
	msgRecordDB := dao.Query().MqttMessageRecord
	list, err := msgRecordDB.Order(msgRecordDB.ID.Desc()).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取消息记录列表").SetData(list)
}

// Save 保存新消息记录
func (s *MqttMessageRecordService) Save(param string) any {
	msgRecordDB := dao.Query().MqttMessageRecord
	var m model.MqttMessageRecord
	if err := json.Unmarshal([]byte(param), &m); err != nil {
		return gf.Failed().SetMsg("添加数据解析失败，" + err.Error())
	}
	//查内容是否存在
	hase, _ := msgRecordDB.Where(msgRecordDB.Payload.Eq(m.Payload)).Count()
	if hase > 0 {
		return gf.Success().SetMsg("数据已经存在")
	}
	if err := msgRecordDB.Create(&m); err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("添加消息记录成功").SetData(m.ID)
}

// Del 删除数据
func (s *MqttMessageRecordService) Del(id int32) any {
	msgRecordDB := dao.Query().MqttMessageRecord
	res, err := msgRecordDB.Where(msgRecordDB.ID.Eq(id)).Delete()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("删除数据成功").SetData(res)
}

// GetContent 获取内容
func (s *MqttMessageRecordService) GetContent(id int32) any {
	msgRecordDB := dao.Query().MqttMessageRecord
	data, err := msgRecordDB.Where(msgRecordDB.ID.Eq(id)).First()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取内容").SetData(data)
}
