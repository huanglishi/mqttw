package gf

import "encoding/json"

// 获取菜单树形
func GetRuleTreeArray(list List, pid int32) List {
	childs := ToolFar(list, pid) //获取pid下的所有数据
	var chridnum List
	for _, v := range childs {
		newdata := GetRuleTreeArray(list, Int32(v["id"]))
		if newdata != nil {
			v["children"] = GetRuleTreeArray(list, Int32(v["id"]))
		}
		chridnum = append(chridnum, v)
	}
	return chridnum
}

// base_tool-获取pid下所有数组
func ToolFar(data List, pid int32) List {
	var mapString List
	for _, v := range data {
		if Int32(v["pid"]) == pid {
			mapString = append(mapString, v)
		}
	}
	return mapString
}

// 获取请求参数id-用于数据保存或更新
func GetEditId(idstr interface{}) int32 {
	id := int32(0)
	if v, ok := idstr.(float64); ok {
		id = int32(v)
	}
	return id
}

// []uint32转 string
func Uint32SliceToString(str []uint32) string {
	if str == nil {
		return ""
	}
	b, err := json.Marshal(str)
	if err != nil {
		return ""
	}
	return string(b)
}

// [][] string 转 string
func DoubleSliceToString(v [][]string) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
