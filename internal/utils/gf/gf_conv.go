package gf

import (
	"strconv"
)

func Int32(v interface{}) int32 {
	switch src := v.(type) {
	case int32:
		return src
	case int:
		return int32(src)
	case int64:
		return int32(src)
	case uint:
		return int32(src)
	case uint32:
		return int32(src)
	case float64:
		return int32(src)
	case string:
		// 字符串数字解析，例如"100"
		num, err := strconv.Atoi(src)
		if err != nil {
			return 0
		}
		return int32(num)
	default:
		return 0
	}
}
