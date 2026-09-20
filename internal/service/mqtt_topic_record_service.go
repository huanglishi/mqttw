package service

import (
	"encoding/json"

	"gofly/internal/dao"
	"gofly/internal/dao/model"
	"gofly/internal/utils/gf"
)

type MqttTopicRecordService struct{}

// GetList 获取消息记录列表
func (g *MqttTopicRecordService) GetList() any {
	topicDB := dao.Query().MqttTopicRecord
	list, err := topicDB.Order(topicDB.ID.Asc()).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取消息记录列表").SetData(list)
}

// Save 保存新消息记录
func (s *MqttTopicRecordService) Save(param string) any {
	topicDB := dao.Query().MqttTopicRecord
	var m model.MqttTopicRecord
	if err := json.Unmarshal([]byte(param), &m); err != nil {
		return gf.Failed().SetMsg("添加数据解析失败，" + err.Error())
	}
	//查内容是否存在
	hase, _ := topicDB.Where(topicDB.Topic.Eq(m.Topic)).Count()
	if hase > 0 {
		topicDB.Where(topicDB.Topic.Eq(m.Topic)).UpdateSimple(topicDB.Qos.Value(m.Qos), topicDB.Retain.Value(m.Retain))
		return gf.Success().SetMsg("数据已经存在")
	}
	if err := topicDB.Create(&m); err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("添加消息记录成功").SetData(m.ID)
}

// Del 删除数据
func (s *MqttTopicRecordService) Del(id int32) any {
	topicDB := dao.Query().MqttTopicRecord
	res, err := topicDB.Where(topicDB.ID.Eq(id)).Delete()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("删除数据成功").SetData(res)
}

// GetContent 获取内容
func (s *MqttTopicRecordService) GetContent(id int32) any {
	topicDB := dao.Query().MqttTopicRecord
	data, err := topicDB.Where(topicDB.ID.Eq(id)).First()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取内容").SetData(data)
}
