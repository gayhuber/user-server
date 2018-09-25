package models

import (
	"encoding/json"
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
	fmt.Println("sendMobileSmsCode:", res)

	resp := baseResp{}
	json.Unmarshal(res.([]byte), &resp)

	if err != nil {
		return err
	}

	if resp.ErrorCode != 0 {
		return errors.New(resp.ErrorMsg)
	}
	return nil
}
