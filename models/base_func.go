package models

import (
	"errors"
	"fmt"
	"log"
	"user-server/tools"
	"user-server/tools/hprose"
)

type args map[string]interface{}

// RespBody 调用 base 方法时返回的结构
type RespBody struct {
	ErrorCode    int                    `json:"errorCode"`
	ErrorMsg     string                 `json:"errorMsg"`
	ResponseData map[string]interface{} `json:"responseData"`
}

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

	hprose.ServiceMp.AddMethod("getServiceProductOrderId", hprose.BaseClient{
		Module: "Order",
		Class:  "Order",
		Func:   "getServiceProductOrderId",
	})

	hprose.ServiceMp.AddMethod("getUserAccoutInfo", hprose.BaseClient{
		Module: "Account",
		Class:  "Account",
		Func:   "getUserAccoutInfo",
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
		// "autoPasswd":  "",
		// "catpcha":     "",
		"ip": ip,
	}

	res, err := hprose.RemoteFunc("sendMobileSmsCode", arg)

	if err != nil {
		return err
	}

	tmp := respHandler(res)

	if tmp["errorCode"] != "0" {
		return errors.New(tmp["errorMsg"].(string))
	}
	return nil
}

func respHandler(res interface{}) (tmp map[string]interface{}) {
	// map 需要初始化一个出来
	tmp = make(map[string]interface{})
	switch res.(type) {
	case nil:
		return tmp
	case map[string]interface{}:
		return res.(map[string]interface{})
	case map[interface{}]interface{}:
		for k, v := range res.(map[interface{}]interface{}) {
			switch k.(type) {
			case string:
				switch v.(type) {
				case map[interface{}]interface{}:
					tmp[k.(string)] = respHandler(v)
					continue
				default:
					tmp[k.(string)] = v
					continue
				}

			default:
				continue
			}
		}
		return tmp
	default:
		// 暂时没遇到更复杂的数据
		log.Println("[respHandler] unknow data:", res)
	}
	return tmp
}

// QuickMobileLogin 手机验证码快速注册/登录
func QuickMobileLogin(mobile, smsCode, countryCode string, sys int) (responseData map[string]interface{}, err error) {
	// return map[string]interface{}{
	// 	"uid":          "20532096",
	// 	"nickname":     "氧气wsc7a",
	// 	"avatar":       "http://img2.soyoung.com/user/5_100_100.png",
	// 	"login_mobile": "18333636949",
	// 	"new_user":     0,
	// }, nil

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
	}

	res, err := hprose.RemoteFunc("quickMobileLogin", arg)
	// fmt.Println("remote QuickMobileLogin:", res, err)

	if err != nil {
		return
	}

	tmp := respHandler(res)
	resp := RespBody{}
	err = tools.Map2Struct(tmp, &resp)

	// 每个接口对成功返回的定义还不一样....
	if resp.ErrorCode == 200 {
		return resp.ResponseData, nil
	}
	err = errors.New(resp.ErrorMsg)
	return
}

// GetSimpleUserInfoByID 获取用户信息
func GetSimpleUserInfoByID(uid int, returnUIDAsKey bool, arrField []string) (info map[string]interface{}, err error) {

	arg := args{
		"uid":            uid,
		"returnUidAsKey": returnUIDAsKey,
		"arrField":       arrField,
	}

	res, err := hprose.RemoteFunc("getSimpleUserInfoById", arg)
	// fmt.Println("remote GetSimpleUserInfoById:", res, err)

	if err != nil {
		return
	}

	tmp := respHandler(res)
	return tmp, nil
}

// GetServiceProductOrderID 获取订单概况
func GetServiceProductOrderID(uid, beging, limit, Type, group int) (info map[string]interface{}, err error) {

	arg := args{
		"uid":   uid,
		"begin": beging,
		"limit": limit,
		"type":  Type,
		"group": group,
	}

	res, err := hprose.RemoteFunc("getServiceProductOrderId", arg)
	// fmt.Println("remote getServiceProductOrderId:", res, err)

	if err != nil {
		return
	}

	info = respHandler(res)
	return
}

// GetUserAccoutInfo 获取用户账户余额
func GetUserAccoutInfo(uid int) (info map[string]interface{}, err error) {
	arg := args{
		"uid": uid,
	}

	// fmt.Println("before exec", arg)
	res, err := hprose.RemoteFunc("getUserAccoutInfo", arg)
	fmt.Println("remote getUserAccoutInfo:", res, err)

	if err != nil {
		return
	}

	info = respHandler(res)
	return
}
