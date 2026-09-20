package service

import (
	"encoding/json"
	"strings"

	"gofly/internal/dao"
	"gofly/internal/dao/model"
	"gofly/internal/dao/query"
	"gofly/internal/utils/gf"
)

type MqttMessageService struct{}

// GetList 获取消息列表
// cid 是connection表数据id，searchTopic是从订阅表获取的主题如：topic/#
func (g *MqttMessageService) GetList(cid int32, searchTopic, msgType string) any {
	msgDB := dao.Query().MqttMessage
	var whereMap []dao.Condition
	whereMap = append(whereMap, msgDB.ConnectionID.Eq(cid))
	if msgType != "all" {
		whereMap = append(whereMap, msgDB.MsgType.Eq(msgType))
	}
	if searchTopic != "all" {
		topicArr := strings.Split(searchTopic, "/")
		if len(topicArr) > 0 {
			whereMap = append(whereMap, msgDB.Topic.Like(topicArr[0]+"%"))
		} else {
			whereMap = append(whereMap, msgDB.Topic.Eq(searchTopic))
		}
	}
	list, err := msgDB.Where(whereMap...).Order(msgDB.ID.Asc()).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取消息列表").SetData(list)
}

// Save 保存/更新消息
func (s *MqttMessageService) Save(param string) any {
	msgDB := dao.Query().MqttMessage
	var m model.MqttMessage
	if err := json.Unmarshal([]byte(param), &m); err != nil {
		return gf.Failed().SetMsg("添加数据解析失败，" + err.Error())
	}
	// 新增
	if m.ID == 0 {
		if err := msgDB.Create(&m); err != nil {
			return gf.Failed().SetMsg(err.Error())
		}
		return gf.Success().SetMsg("添加消息成功").SetData(m.ID)
	}
	res, err := msgDB.Where(msgDB.ID.Eq(m.ID)).Updates(m)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("更新消息成功").SetData(res.RowsAffected)
}

// Del 删除数据
func (s *MqttMessageService) Del(id int32) any {
	msgDB := dao.Query().MqttMessage
	res, err := msgDB.Where(msgDB.ID.Eq(id)).Delete()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("删除数据成功").SetData(res)
}

// ClearMessage 清除指定连接数据
func (s *MqttMessageService) ClearMessage(id int32) any {
	msgDB := dao.Query().MqttMessage
	res, err := msgDB.Where(msgDB.ConnectionID.Eq(id)).Delete()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("清除连接数据成功").SetData(res)
}

// GetContent 获取内容
func (s *MqttMessageService) GetContent(id int32) any {
	msgDB := dao.Query().MqttMessage
	data, err := msgDB.Where(msgDB.ID.Eq(id)).First()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取内容").SetData(data)
}

// ClearRecord 清除主题和内容记录
func (s *MqttMessageService) ClearRecord() any {
	q := dao.Query()
	err := q.Transaction(func(tx *query.Query) error {
		db := tx.UnderlyingDB()
		// 1. mqtt_message_record
		if err := db.Exec("DELETE FROM mqtt_message_record").Error; err != nil {
			return err
		}
		if err := db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = ?", "mqtt_message_record").Error; err != nil {
			return err
		}
		// 2. mqtt_topic_record
		if err := db.Exec("DELETE FROM mqtt_topic_record").Error; err != nil {
			return err
		}
		if err := db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = ?", "mqtt_topic_record").Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return gf.Failed().SetMsg("清空记录失败：" + err.Error())
	}
	return gf.Success().SetMsg("全部记录已清空")
}
