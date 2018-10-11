package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
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

	userInfo := tokenBody{}
	json.Unmarshal(resp, &userInfo)
	uid, _ := tools.String2Int(userInfo.UID)

	info, err := GetSimpleUserInfoByID(uid, false, []string{"uid", "user_name", "avatar", "login_mobile"})
	if err != nil {
		return 400, err
	}
	mobile, ok := info["login_mobile"].(string)
	if ok {
		info["login_mobile"] = mobile[:3] + "*****" + mobile[8:]
	} else {
		info["login_mobile"] = ""
	}

	orderInfo, err := GetServiceProductOrderID(uid, 1, 10, 1, 1)
	if err != nil {
		return 400, err
	}

	_, ok = orderInfo["total"]
	if ok {
		info["unpaid_order"] = orderInfo["total"]
	} else {
		info["unpaid_order"] = 0
	}

	info["kefu_mobile"] = "4001816660"

	return 200, info
}

func (auth *MobileAuth) homeNew() (code int, obj interface{}) {
	// 获取用户 uid
	userInfo, err := NewSession(auth.Token, "mobile").infoStruct()
	if err != nil {
		return 400, err
	}
	uid, _ := tools.String2Int(userInfo.UID)

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)
	defer close(errChan)

	// 获取用户信息
	infoChan := make(chan map[string]interface{}, 1)
	defer close(infoChan)
	go func(infoChan chan<- map[string]interface{}, errChan chan<- error) {
		defer wg.Done()
		info, err := GetSimpleUserInfoByID(uid, false, []string{"uid", "user_name", "avatar", "login_mobile"})

		if err != nil {
			errChan <- err
			return
		}
		if len(info) == 0 {
			errChan <- errors.New("user not found")
		}
		mobile, ok := info["login_mobile"].(string)

		if ok {
			info["login_mobile"] = mobile[:3] + "*****" + mobile[8:]
		} else {
			info["login_mobile"] = ""
		}

		infoChan <- info
	}(infoChan, errChan)

	// 获取订单信息
	accountTotalChan := make(chan interface{}, 1)
	defer close(accountTotalChan)
	go func(accountTotalChan chan<- interface{}, errChan chan<- error) {
		defer wg.Done()
		orderInfo, err := GetServiceProductOrderID(uid, 1, 10, 1, 1)
		fmt.Println("GetServiceProductOrderID result:", orderInfo, err)
		if err != nil {
			errChan <- err
			return
		}
		_, ok := orderInfo["total"]
		if ok {
			accountTotalChan <- orderInfo["total"]
		} else {
			accountTotalChan <- 0
		}
		return
	}(accountTotalChan, errChan)

	wg.Wait()

	if len(errChan) > 0 {
		return 400, <-errChan
	}

	info := <-infoChan
	info["kefu_mobile"] = "4001816660"
	info["unpaid_order"] = <-accountTotalChan

	return 200, info
}

func (auth *MobileAuth) account() (code int, obj interface{}) {
	userInfo, err := NewSession(auth.Token, "mobile").infoStruct()
	if err != nil {
		return 400, err
	}

	uid, _ := tools.String2Int(userInfo.UID)

	account, err := GetUserAccoutInfo(uid)

	if err != nil {
		return 400, err
	}

	var balance string
	if len(account) > 0 {
		balance = account["total_amount"].(string)
	} else {
		balance = "0.00"
	}

	return 200, lib.H{
		"balance": balance,
	}
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
func MobileHome(session *lib.Session) {
	auth := MobileAuth{}
	auth.setParams(session.Request.Params)
	code, resp := auth.homeNew()
	session.Send(code, resp)
}

// MobileAccount 钱包信息
func MobileAccount(session *lib.Session) {
	auth := MobileAuth{}
	auth.setParams(session.Request.Params)
	code, resp := auth.account()
	session.Send(code, resp)
}
