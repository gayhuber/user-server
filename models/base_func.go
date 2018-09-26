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

	// 这个变相转换一下输出, 用 json.Unmarshal 会出现map格式转换问题
	tmp := map[string]string{}
	for k, v := range res.(map[interface{}]interface{}) {
		key := fmt.Sprintf("%v", k)
		value := fmt.Sprintf("%v", v)
		tmp[key] = value
	}

	if err != nil {
		return err
	}

	if tmp["errorCode"] != "0" {
		return errors.New(tmp["errorMsg"])
	}
	return nil
}
