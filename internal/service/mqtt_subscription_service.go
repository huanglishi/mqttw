package service

import (
	"encoding/json"

	"gofly/internal/dao"
	"gofly/internal/dao/model"
	"gofly/internal/utils/gf"
)

type MqttSubscriptionService struct{}

// GetList 获取订阅列表
// cid 是connection表数据id
func (g *MqttSubscriptionService) GetList(cid int32) any {
	subDB := dao.Query().MqttSubscription
	list, err := subDB.Where(subDB.ConnectionID.Eq(cid)).Order(subDB.ID.Asc()).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取订阅列表").SetData(list)
}

// Save 保存/更新订阅数据
func (s *MqttSubscriptionService) Save(param string) any {
	subDB := dao.Query().MqttSubscription
	var m model.MqttSubscription
	if err := json.Unmarshal([]byte(param), &m); err != nil {
		return gf.Failed().SetMsg("添加表单数据解析失败，" + err.Error())
	}
	if m.Topic == "" {
		return gf.Failed().SetMsg("主题不能为空")
	}
	// 重复校验：同一连接下 Topic 不允许重复（更新时排除自身）
	if m.ConnectionID > 0 && m.Topic != "" {
		var count int64
		var err error
		if m.ID > 0 {
			count, err = subDB.Where(subDB.ConnectionID.Eq(m.ConnectionID), subDB.Topic.Eq(m.Topic), subDB.ID.Neq(m.ID)).Count()
		} else {
			count, err = subDB.Where(subDB.ConnectionID.Eq(m.ConnectionID), subDB.Topic.Eq(m.Topic)).Count()
		}
		if err != nil {
			return gf.Failed().SetMsg(err.Error())
		}
		if count > 0 {
			return gf.Failed().SetMsg("该连接下已存在相同主题的订阅：" + m.Topic)
		}
	}
	// 新增
	if m.ID == 0 {
		if err := subDB.Create(&m); err != nil {
			return gf.Failed().SetMsg(err.Error())
		}
		return gf.Success().SetMsg("添加订阅成功").SetData(m.ID).SetExdata("add")
	}
	var oldTopic string
	subDB.Where(subDB.ID.Eq(m.ID)).Select(subDB.Topic).Scan(&oldTopic)
	_, err := subDB.Where(subDB.ID.Eq(m.ID)).Updates(m)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("更新订阅成功").SetData(m.ConnectionID).SetExdata(oldTopic)
}

// Disable 禁用
// status 禁用状态:0=正常,1=禁用
func (s *MqttSubscriptionService) Disable(id, status int32) any {
	subDB := dao.Query().MqttSubscription
	res, err := subDB.Where(subDB.ID.Eq(id)).Update(subDB.Disabled, status)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("状态变更成功").SetData(res)
}

// Del 删除数据
func (s *MqttSubscriptionService) Del(id int32) any {
	subDB := dao.Query().MqttSubscription
	res, err := subDB.Where(subDB.ID.Eq(id)).Delete()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("删除数据成功").SetData(res)
}

// GetContent 获取内容
func (s *MqttSubscriptionService) GetContent(id int32) any {
	subDB := dao.Query().MqttSubscription
	data, err := subDB.Where(subDB.ID.Eq(id)).First()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取内容").SetData(data)
}
