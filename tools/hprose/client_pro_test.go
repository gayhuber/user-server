package hprose

import (
	"encoding/json"
	"testing"
)

// 测试新的远程调用方式, 因为用到了根据服务名调用不同的base server
// 同时使用了hprose/io组件, 已解决出现的 map[interface{}]interface{} 的问题
func TestRemoteFuncPro1(t *testing.T) {

	args := []interface{}{
		20532096,
		false,
		[]string{"uid", "user_name", "avatar", "login_mobile"},
	}

	info, err := RemoteFuncPro("php_base_server", "User", "UserNew", "getSimpleUserInfoById", args)
	if err != nil {
		t.Error(err)
		return
	}

	jsonByte, err := json.Marshal(info)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("result: ", info, err, string(jsonByte))
}

func TestRemoteFuncPro2(t *testing.T) {
	args := []interface{}{
		18333636949,
	}

	info, err := RemoteFuncPro("php_base_server", "System", "XyToken", "getXyOpenKey", args)
	if err != nil {
		t.Error(err)
		return
	}

	jsonByte, err := json.Marshal(info)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("result: ", info, err, string(jsonByte))
}
