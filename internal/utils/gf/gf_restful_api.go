package gf

import (
	"time"
)

// 组装Restful API 接口返回数据
// 返回信息主体
type (
	R struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Exdata  interface{} `json:"exdata"`
		Time    int64       `json:"time"`
	}
)

// 错误码
var (
	succCode = 0 // 成功
	errCode  = 1 // 失败
)

// 设置返回内容
func (r *R) SetData(data interface{}) *R {
	r.Data = data
	return r
}

// 设置返回扩展内容内容
func (r *R) SetExdata(exdata interface{}) *R {
	r.Exdata = exdata
	return r
}

// 设置编码
func (r *R) SetCode(code int) *R {
	r.Code = code
	return r
}

// 设置返回提示信息
func (r *R) SetMsg(msg string) *R {
	r.Message = msg
	return r
}

// 返回成功内容
func Success() *R {
	r := &R{}
	r.Message = "Success"
	r.Code = succCode
	r.Time = time.Now().Unix()
	return r
}

// 返回失败内容
func Failed() *R {
	r := &R{}
	r.Message = "Fail"
	r.Code = errCode
	r.Data = false
	r.Time = time.Now().Unix()
	return r
}
