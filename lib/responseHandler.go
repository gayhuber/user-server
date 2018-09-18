package lib

import (
	"fmt"
	"sync"
)

type statusMap struct {
	RLock sync.RWMutex
	Body  map[int]string
}

// 这个文件是根据 code 编码来标记错误的
var statusmp *statusMap

func init() {
	statusmp = &statusMap{}

	statusmp.Body = map[int]string{
		500:  "程序错误",
		200:  "成功",
		1002: "参数错误",
		1004: "禁止访问",
		1005: "该服务暂时不可用",
		2001: "数据保存不成功",
		2002: "无效的 open_id",
		2003: "无效的 token",
		2004: "注册失败",
		2005: "渠道来源不合法",
		4001: "需要登陆",
		4005: "用户不存在",
	}
}

// GetCodeMsg 将 code 转换成文字
func GetCodeMsg(code int) (msg string) {
	msg, _ = statusmp.Body[code]
	return
}

// ResponseHandler 处理返回格式
// 保证返回内容中
func ResponseHandler(code int, obj interface{}) Response {
	var resp Response

	resp.Code = code
	resp.Message = GetCodeMsg(code)
	resp.Data = obj

	if code != 200 {
		switch obj.(type) {
		case H:
			data := obj.(H)
			msg, ok := data["msg"]
			if ok {
				resp.Data = msg.(string)
			} else {
				resp.Data = ""
			}
		default:
			resp.Message = fmt.Sprint(obj)
			resp.Data = ""
		}

	}
	return resp
}

// TCPError 统一的错误处理(待定)
func TCPError(code int, extra ...interface{}) {
	if len(extra) > 0 {

	}
}

// ErrorFormat 错误格式(待定)
type ErrorFormat struct {
	Code  int
	Extra string
}
