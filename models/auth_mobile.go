package models

import (
	"fmt"
	"strconv"
	"time"
	"user-server/lib"
	"user-server/tools"
)

// USERKEY 加密通用 salt
const USERKEY = "ncisnotnaocan!@#2018"

// MobileAuth 来自 mobile 渠道的参数格式
type MobileAuth struct {
	Type        string
	Mobile      string
	Code        string
	CountryCode string
	Sys         int
	Ext         map[string]interface{}
	IP          string
	Token       string
}

func (auth *MobileAuth) getName() string {
	return "MobileAuth"
}

func (auth *MobileAuth) register() (code int, obj interface{}) {
	resp, err := QuickMobileLogin(auth.Mobile, auth.Code, auth.CountryCode, 11)
	if err != nil {
		return 400, err
	}

	uid, _ := strconv.Atoi(resp["uid"].(string))
	token := auth.generateToken(uid)
	secKey := auth.generateSecKey(uid)

	resp["token"] = token
	resp["sec_key"] = secKey

	// 将用户登录信息存到缓存中一份
	NewSession(token, "mobile").store(resp)

	return 200, lib.H{
		"sec_key":  secKey,
		"token":    token,
		"lifetime": tools.CalculateTTL(),
	}
}

func (auth *MobileAuth) login() (code int, obj interface{}) {
	return 401, nil
}
func (auth *MobileAuth) info() (code int, obj interface{}) {
	resp, err := NewSession(auth.Token, "mobile").info()
	if err != nil {
		return 400, err
	}
	return 200, resp
}
func (auth *MobileAuth) home() (code int, obj interface{}) {
	resp, err := NewSession(auth.Token, "mobile").info()
	if err != nil {
		return 400, err
	}
	return 200, resp
}
func (auth *MobileAuth) generateToken(uid int) string {
	tmp := fmt.Sprintf("%d%s%s", uid, USERKEY, time.Now().String())
	return tools.MD5String([]byte(tmp))
}
func (auth *MobileAuth) generateSecKey(uid int) string {
	tmp := fmt.Sprintf("%d%s%s", uid, USERKEY, time.Now().String())
	return tools.MD5String([]byte(tmp))
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

	if Token, ok := params["token"]; ok {
		auth.Token = Token.(string)
	}

	// 平台编号
	if sys, ok := params["sys"]; ok {
		auth.Sys = int(sys.(float64))
	}

	// 短信验证码
	if code, ok := params["code"]; ok {
		auth.Code = code.(string)
	}
	auth.Ext = params

	return
}

// sms 调用 base 的发送短信接口
func (auth *MobileAuth) sms() (code int, obj interface{}) {

	fmt.Println("sms:", auth)
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

// MobileHome 首页信息
// TODO: 改成和php 版一样的返回内容
func MobileHome(session *lib.Session) {
	auth := MobileAuth{}
	auth.setParams(session.Request.Params)
	code, resp := auth.home()
	session.Send(code, resp)
}
