package models

import (
	"errors"
	"fmt"
	"user-server/tools/hprose"
)

type args map[string]interface{}

func init() {
	// 注册方法
	hprose.ServiceMp.AddMethod("getXyOpenKey", hprose.BaseClient{
		Module: "System",
		Class:  "XyToken",
		Func:   "getXyOpenKey",
	})

	hprose.ServiceMp.AddMethod("sendMobileSmsCode", hprose.BaseClient{
		Module: "Passport",
		Class:  "Core\\MobileSmsCode",
		Func:   "sendMobileSmsCode",
	})

	hprose.ServiceMp.AddMethod("quickMobileLogin", hprose.BaseClient{
		Module: "Passport",
		Class:  "Core\\Login",
		Func:   "quickMobileLogin",
	})

	hprose.ServiceMp.AddMethod("getSimpleUserInfoById", hprose.BaseClient{
		Module: "User",
		Class:  "UserNew",
		Func:   "getSimpleUserInfoById",
	})
}

// GetXyOpenKey 获取 openkey
func GetXyOpenKey(mobile string) (string, error) {
	res, err := hprose.RemoteFunc("getXyOpenKey", args{
		"mobile": mobile,
	})

	return res.(string), err
}

type baseResp struct {
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

// SendMobileSmsCode 发送远程短信
func SendMobileSmsCode(mobile, countryCode, ip, Type string) error {
	openKey, err := GetXyOpenKey(mobile)
	if err != nil {
		return err
	}

	arg := args{
		"mobile":      mobile,
		"countryCode": countryCode,
		"type":        Type,
		"key":         openKey,
		"autoPasswd":  "",
		"catpcha":     "",
		"ip":          ip,
	}

	res, err := hprose.RemoteFunc("sendMobileSmsCode", arg)

	if err != nil {
		return err
	}

	// 这个变相转换一下输出, 用 json.Unmarshal 会出现map格式转换问题
	// tmp := map[string]string{}
	// for k, v := range res.(map[interface{}]interface{}) {
	// 	key := fmt.Sprintf("%v", k)
	// 	value := fmt.Sprintf("%v", v)
	// 	tmp[key] = value
	// }
	tmp := respHandler(res)

	if tmp["errorCode"] != "0" {
		return errors.New(tmp["errorMsg"])
	}
	return nil
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
		"lver":    "",
		"version": "",
	}

	res, err := hprose.RemoteFunc("quickMobileLogin", arg)
	fmt.Println("remote QuickMobileLogin:", res, err)

	if err != nil {
		return
	}

	tmp := respHandler(res)

	if tmp["errorCode"] != "0" {
		err = errors.New(tmp["errorMsg"])
		return
	}
	return tmp["responseData"], nil
}

// GetSimpleUserInfoByID 获取用户信息
func GetSimpleUserInfoByID(uid int, returnUIDAsKey bool, arrField []string) (info map[string]string, err error) {

	arg := args{
		"uid":            uid,
		"returnUidAsKey": returnUIDAsKey,
		"arrField":       arrField,
	}

	res, err := hprose.RemoteFunc("getSimpleUserInfoById", arg)
	fmt.Println("remote GetSimpleUserInfoById:", res, err)

	if err != nil {
		return
	}

	tmp := respHandler(res)
	fmt.Println("this is tmp:", tmp)
	return tmp, nil
}
