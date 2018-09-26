package models

import (
	"fmt"
	"user-server/lib"
)

// MobileAuth 来自 mobile 渠道的参数格式
type MobileAuth struct {
	Type        string
	Mobile      string
	Code        string
	CountryCode string
	Sys         int
	Ext         map[string]interface{}
	IP          string
}

func (auth *MobileAuth) getName() string {
	return "MobileAuth"
}

func (auth *MobileAuth) register() (code int, obj interface{}) {
	return 200, lib.H{
		"result": "mobile register",
	}
}

func (auth *MobileAuth) login() (code int, obj interface{}) {
	return 401, nil
}
func (auth *MobileAuth) info() (code int, obj interface{}) {
	return 401, nil
}
func (auth *MobileAuth) home() (code int, obj interface{}) {
	return 401, nil
}
func (auth *MobileAuth) setParams(params map[string]interface{}) {
	if mobile, ok := params["mobile"]; ok {
		auth.Mobile = mobile.(string)
	}
	if countryCode, ok := params["country_code"]; ok {
		auth.CountryCode = countryCode.(string)
	}
	if IP, ok := params["ip"]; ok {
		auth.IP = IP.(string)
	}
	auth.Ext = params

	return
}

// sms 调用 base 的发送短信接口
func (auth *MobileAuth) sms() (code int, obj interface{}) {

	fmt.Println(auth)
	// 8 代表着快速登录
	err := SendMobileSmsCode(auth.Mobile, auth.CountryCode, auth.IP, "8")

	if err != nil {
		return 400, err
	}

	return 200, nil
}

// 以下是 mobile 独有的接口:

// MobileSms 发送短信验证码
func MobileSms(session *lib.Session) {
	auth := MobileAuth{}
	auth.setParams(session.Request.Params)
	code, resp := auth.sms()
	session.Send(code, resp)
}