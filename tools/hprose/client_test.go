package hprose

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type args map[string]interface{}

func respHandler(res interface{}) map[string]string {
	tmp := map[string]string{}
	for k, v := range res.(map[interface{}]interface{}) {
		key := fmt.Sprintf("%v", k)
		value := fmt.Sprintf("%v", v)
		tmp[key] = value
	}
	return tmp
}

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

func TestCallback3(t *testing.T) {
	ServiceMp.AddMethod("getXyOpenKey", BaseClient{
		Module: "System",
		Class:  "XyToken",
		Func:   "getXyOpenKey",
	})
	ServiceMp.AddMethod("quickMobileLogin", BaseClient{
		Module: "Passport",
		Class:  "Core\\Login",
		Func:   "quickMobileLogin",
	})

	resp, err := QuickMobileLogin("18500215317", "482076", "0086", 11)

	if err != nil {
		t.Error(err)
	}

	fmt.Println("result:", resp)

	tmp := map[string]interface{}{}
	json.Unmarshal([]byte(resp), &tmp)
	fmt.Println("tmp output:", tmp)
}

// GetXyOpenKey 获取 openkey
func GetXyOpenKey(mobile string) (string, error) {
	res, err := RemoteFunc("getXyOpenKey", args{
		"mobile": mobile,
	})

	fmt.Println("remote openKey:", res, err)

	return res.(string), err
}

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

// QuickMobileLogin 手机验证码快速注册/登录
func QuickMobileLogin(mobile, smsCode, countryCode string, sys int) (responseData string, err error) {
	openKey, err := GetXyOpenKey(mobile)
	if err != nil {
		return
	}

	arg := args{
		"mobile":      mobile,
		"smsCode":     smsCode,
		"countryCode": countryCode,
		"key":         openKey,
		"extInfo": args{
			"sys": sys,
		},
		// "lver":    "",
		// "version": "",
	}

	res, err := RemoteFunc("quickMobileLogin", arg)
	fmt.Println("remote QuickMobileLogin:", res, err)

	if err != nil {
		return
	}

	// TODO: respHandler 需要把responseData转成[]byte 类型的
	tmp := respHandler(res)

	// 每个接口对成功返回的定义还不一样....
	if tmp["errorCode"] != "200" {
		err = errors.New(tmp["errorMsg"])
		return
	}
	return tmp["responseData"], nil
}
