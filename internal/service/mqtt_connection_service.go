package service

import (
	"encoding/json"
	"strings"
	"time"

	"gofly/internal/dao"
	"gofly/internal/dao/model"
	"gofly/internal/dao/query"
	"gofly/internal/utils/gf"

	"gorm.io/gorm"
)

type MqttConnectionService struct{}

// 启动服务时还原数据库连接状态
func init() {
	go func() {
		//运用启动时延迟执行
		time.Sleep(time.Second * 3)
		connectionDB := dao.Query().MqttConnection
		connectionDB.Where(connectionDB.Connected.Eq(1)).Update(connectionDB.Connected, 0)
	}()
}

// camelToSnake 前端键名 clientId -> client_id，与表列名对齐
func camelToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// normalizeData 更新用：键名转 snake_case、数字转 int32、嵌套对象序列化为 JSON 文本
func normalizeData(data map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(data))
	for k, v := range data {
		col := camelToSnake(k)
		switch col {
		case "id", "create_time", "update_time":
			continue
		}
		switch t := v.(type) {
		case float64:
			if t == float64(int32(t)) {
				v = int32(t)
			}
		case map[string]interface{}, []interface{}:
			if b, err := json.Marshal(t); err == nil {
				v = string(b)
			}
		}
		out[col] = v
	}
	return out
}

// GetList 获取连接/分组树
func (g *MqttConnectionService) GetList() any {
	connectionDB := dao.Query().MqttConnection
	var list gf.List
	err := connectionDB.Select(
		connectionDB.ID,
		connectionDB.Pid,
		connectionDB.Title,
		connectionDB.IsGroup,
		connectionDB.Weigh,
	).Order(connectionDB.IsGroup.Desc(), connectionDB.Weigh.Asc()).Scan(&list)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	if len(list) > 0 {
		list = gf.GetRuleTreeArray(list, 0)
	}
	return gf.Success().SetMsg("获取数据成功").SetData(list)
}

// SaveGroup 保存/更新分组数据
func (s *MqttConnectionService) SaveGroup(param string) any {
	connectionDB := dao.Query().MqttConnection
	var m model.MqttConnection
	if err := json.Unmarshal([]byte(param), &m); err != nil {
		return gf.Failed().SetMsg("添加表单数据解析失败，" + err.Error())
	}
	// 新增
	if m.ID == 0 {
		if err := connectionDB.Create(&m); err != nil {
			return gf.Failed().SetMsg(err.Error())
		}
		if _, err := connectionDB.Where(connectionDB.ID.Eq(m.ID)).Update(connectionDB.Weigh, m.ID); err != nil {
			return gf.Failed().SetMsg(err.Error())
		}
		return gf.Success().SetMsg("添加分组成功").SetData(m.ID)
	}
	res, err := connectionDB.Where(connectionDB.ID.Eq(m.ID)).Updates(m)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("更新分组成功").SetData(res.RowsAffected)
}

// SaveData 保存/更新数据
func (s *MqttConnectionService) SaveData(param string) any {
	connectionDB := dao.Query().MqttConnection
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(param), &data); err != nil {
		return gf.Failed().SetMsg("添加表单数据解析失败，" + err.Error())
	}
	md := normalizeData(data)
	var f_id = gf.GetEditId(data["id"])
	// 新增
	if f_id == 0 {
		md["created_at"] = time.Now()
		if tx := connectionDB.UnderlyingDB().Model(&model.MqttConnection{}).Create(&md); tx.Error != nil {
			return gf.Failed().SetMsg(tx.Error.Error())
		}
		if _, err := connectionDB.Where(connectionDB.ID.Eq(gf.Int32(md["id"]))).Update(connectionDB.Weigh, md["id"]); err != nil {
			return gf.Failed().SetMsg(err.Error())
		}
		return gf.Success().SetMsg("添加数据成功").SetData(md["id"])
	}
	// 更新
	if len(md) == 0 {
		return gf.Failed().SetMsg("没有可更新的字段")
	}
	_, err := connectionDB.Where(connectionDB.ID.Eq(int32(f_id))).Updates(md)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("更新数据成功").SetData(f_id)
}

// Del 删除数据
func (s *MqttConnectionService) Del(id int32) any {
	connectionDB := dao.Query().MqttConnection
	var is_group int32
	connectionDB.Where(connectionDB.ID.Eq(id)).Select(connectionDB.IsGroup).Scan(&is_group)
	_, err := connectionDB.Where(connectionDB.ID.Eq(id)).Delete()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	//删除对应订阅主题
	subDB := dao.Query().MqttSubscription
	subDB.Where(subDB.ConnectionID.Eq(id)).Delete()
	msgDB := dao.Query().MqttMessage
	msgDB.Where(msgDB.ConnectionID.Eq(id)).Delete()
	connectedIds := make([]int32, 0) // 全部在连接中id
	//删除分组下的所有数据（递归删除多级子分组及其下连接）
	if is_group > 0 {
		// 递归收集该分组下全部后代 ID
		childIds := make([]int32, 0) // 全部后代（含子分组、连接）
		connIds := make([]int32, 0)  // 后代中的连接（需级联删订阅、消息）
		if err := s.collectChildIds([]int32{id}, &childIds, &connIds); err != nil {
			return gf.Failed().SetMsg("查询分组子数据失败，" + err.Error())
		}
		if len(childIds) > 0 {
			// 先删后代连接的订阅主题、消息记录
			if len(connIds) > 0 {
				subDB.Where(subDB.ConnectionID.In(connIds...)).Delete()
				msgDB.Where(msgDB.ConnectionID.In(connIds...)).Delete()
			}
			//获取连接中的数据id
			connectionDB.Where(connectionDB.ID.In(childIds...), connectionDB.Connected.Eq(1)).Pluck(connectionDB.ID, &connectedIds)
			// 再删除全部子分组与连接
			if _, err := connectionDB.Where(connectionDB.ID.In(childIds...)).Delete(); err != nil {
				return gf.Failed().SetMsg(err.Error())
			}
		}
	}
	return gf.Success().SetMsg("删除数据成功").SetData(connectedIds)
}

// collectChildIds 递归收集 parentIds 下的全部后代 ID（支持多级嵌套分组）
// childIds：收集到的全部后代（子分组 + 连接）；connIds：其中 is_group=0 的连接 ID
func (s *MqttConnectionService) collectChildIds(parentIds []int32, childIds *[]int32, connIds *[]int32) error {
	// 递归终止：没有待查的父级
	if len(parentIds) == 0 {
		return nil
	}
	connectionDB := dao.Query().MqttConnection
	children, err := connectionDB.Where(connectionDB.Pid.In(parentIds...)).Find()
	if err != nil {
		return err
	}
	nextParentIds := make([]int32, 0, len(children))
	for _, c := range children {
		*childIds = append(*childIds, c.ID)
		nextParentIds = append(nextParentIds, c.ID) // 继续向下递归
		if c.IsGroup == 0 {
			*connIds = append(*connIds, c.ID)
		}
	}
	// 以当前层 ID 作为下一层的父级继续查找
	return s.collectChildIds(nextParentIds, childIds, connIds)
}

// Update 更新数据
func (s *MqttConnectionService) Update(param string) any {
	connectionDB := dao.Query().MqttConnection
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(param), &data); err != nil {
		return gf.Failed().SetMsg("更新表单数据解析失败，" + err.Error())
	}
	res, err := connectionDB.Where(connectionDB.ID.Eq(gf.Int32(data["id"]))).Updates(data)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("更新数据成功").SetData(res)
}

// GetContent 获取内容
func (s *MqttConnectionService) GetContent(id int32) any {
	connectionDB := dao.Query().MqttConnection
	data, err := connectionDB.Where(connectionDB.ID.Eq(id)).First()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取内容").SetData(data)
}

// GetInfo 获取指定字段内容
func (s *MqttConnectionService) GetInfo(id int32) any {
	connectionDB := dao.Query().MqttConnection
	var data map[string]any
	err := connectionDB.Where(connectionDB.ID.Eq(id)).Select(
		connectionDB.ID,
		connectionDB.Title,
		connectionDB.ClientID,
		connectionDB.Username,
		connectionDB.Password,
		connectionDB.Keepalive,
		connectionDB.Clean,
		connectionDB.MqttVersion,
		connectionDB.Connected,
	).Scan(&data)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("获取指定字段内容").SetData(data)
}

// Copy 复制连接数据
func (s *MqttConnectionService) Copy(id int32) any {
	connectionDB := dao.Query().MqttConnection
	copyData, err := connectionDB.Where(connectionDB.ID.Eq(id)).First()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	copyData.Title = copyData.Title + "_copy"
	copyData.ClientID = copyData.ClientID + "_copy"
	copyData.Connected = 0
	err = connectionDB.Omit(connectionDB.ID, connectionDB.CreatedAt, connectionDB.UpdatedAt).Create(copyData)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetMsg("复制连接数据成功").SetData(copyData.ID)
}

// UpdateOrder 批量更新排序：pid（父级归属）+ weigh（同级顺序）
// param 为 JSON 数组字符串：[{"id":1,"pid":0,"weigh":1},{"id":2,"pid":1,"weigh":1},...]
// 前端拖拽结束后把整棵树扁平化传入（同层兄弟按位置编号 weigh=1,2,3...，根级 pid=0）
func (s *MqttConnectionService) UpdateOrder(param string) any {
	var list []model.MqttConnection
	if err := json.Unmarshal([]byte(param), &list); err != nil {
		return gf.Failed().SetMsg("排序数据解析失败，" + err.Error())
	}
	if len(list) == 0 {
		return gf.Failed().SetMsg("没有可排序的数据")
	}
	// 事务批量更新，保证 pid + weigh 要么全部成功、要么全部回滚
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		for _, item := range list {
			if err := tx.Model(&model.MqttConnection{}).
				Where("id = ?", item.ID).
				Updates(map[string]interface{}{
					"pid":   item.Pid,
					"weigh": item.Weigh,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return gf.Failed().SetMsg("排序更新失败，" + err.Error())
	}
	return gf.Success().SetMsg("排序更新成功").SetData(len(list))
}

// 清除全部数据
func (s *MqttConnectionService) ClearData() any {
	q := dao.Query()
	err := q.Transaction(func(tx *query.Query) error {
		db := tx.UnderlyingDB()
		// 1. mqtt_connection
		if err := db.Exec("DELETE FROM mqtt_connection").Error; err != nil {
			return err
		}
		if err := db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = ?", "mqtt_connection").Error; err != nil {
			return err
		}
		// 2. mqtt_subscription
		if err := db.Exec("DELETE FROM mqtt_subscription").Error; err != nil {
			return err
		}
		if err := db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = ?", "mqtt_subscription").Error; err != nil {
			return err
		}
		// 3. mqtt_message
		if err := db.Exec("DELETE FROM mqtt_message").Error; err != nil {
			return err
		}
		if err := db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = ?", "mqtt_message").Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return gf.Failed().SetMsg("清空数据失败：" + err.Error())
	}
	return gf.Success().SetMsg("全部数据已清空")
}
