package hprose

import (
	"testing"
)

func TestCallback(t *testing.T) {
	args := map[string]interface{}{
		"mobile": "18333636949",
	}

	// 注册方法
	ServiceMp.AddMethod("XyToken", BaseClient{
		Module: "System",
		Class:  "XyToken",
		Func:   "getXyOpenKey",
	})

	// 调用方法, 并传递参数进去, 获取远程 base 的调用结果
	res, err := RemoteFunc("XyToken", args)

	t.Log(res, err)
}
