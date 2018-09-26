package hprose

import (
	// "encoding/json"
	// "errors"
	"fmt"
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

func TestCallback2(t *testing.T) {
	ServiceMp.AddMethod("getSimpleUserInfoById", BaseClient{
		Module: "User",
		Class:  "UserNew",
		Func:   "getSimpleUserInfoById",
	})

	field := []string{"uid", "user_name", "avatar", "login_mobile"}
	resp, err := GetSimpleUserInfoById(20532239, false, field)

	if err != nil {
		t.Error(err)
	}

	fmt.Println("result:", resp)
}

type args map[string]interface{}

// GetSimpleUserInfoById 获取用户信息
func GetSimpleUserInfoById(uid int, returnUidAsKey bool, arrField []string) (info map[string]string, err error) {

	arg := args{
		"uid":            uid,
		"returnUidAsKey": returnUidAsKey,
		"arrField":       arrField,
	}

	res, err := RemoteFunc("getSimpleUserInfoById", arg)
	fmt.Println("remote GetSimpleUserInfoById:", res, err)

	if err != nil {
		return
	}

	tmp := respHandler(res)
	fmt.Println("this is tmp:", tmp)
	return tmp, nil
}

func respHandler(res interface{}) map[string]string {
	tmp := map[string]string{}
	for k, v := range res.(map[interface{}]interface{}) {
		key := fmt.Sprintf("%v", k)
		value := fmt.Sprintf("%v", v)
		tmp[key] = value
	}
	return tmp
}
