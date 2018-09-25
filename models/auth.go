// Package models 这里写业务逻辑
package models

import (
	"github.com/pkg/errors"
	"user-server/lib"
)

type authHandler interface {
	getName() string
	register() (code int, obj interface{})
	info() (code int, obj interface{})
	login() (code int, obj interface{})
	home() (code int, obj interface{})
	setParams(params map[string]interface{})
}

var handlerMap map[string]authHandler

// transfer 负责转换渠道和给参数集合赋值
func transfer(tp string, params map[string]interface{}) (hd authHandler, err error) {
	// 根据 type 来加载不同的实体
	switch tp {
	case "h5":
		hd = &H5Auth{}
	}
	if hd == nil {
		err = errors.New("not found handler")
		return
	}
	hd.setParams(params)
	return
}

// AuthRegister 注册用户
func AuthRegister(session *lib.Session) {
	hd, err := transfer(session.Request.Params["type"].(string), session.Request.Params)
	if err != nil {
		session.Send(500, err)
	}
	code, resp := hd.register()

	session.Log.Info(resp, "RESPONSE")
	session.Send(code, resp)
}

// AuthLogin 用户登录
func AuthLogin(session *lib.Session) {
	reqType := session.Request.Params["type"].(string)
	hd, err := transfer(reqType, session.Request.Params)
	if err != nil {
		session.Send(500, err)
	}
	code, resp := hd.login()

	// 将用户登录信息存到缓存中一份
	data := resp.(lib.H)
	NewSession(data["token"].(string), reqType).store(data)

	session.Log.Info(resp, "RESPONSE")
	session.Send(code, resp)
}

func AuthInfo(session *lib.Session) {
	reqType := session.Request.Params["type"].(string)
	hd, err := transfer(reqType, session.Request.Params)
	if err != nil {
		session.Send(500, err)
	}

	code, resp := hd.info()

	session.Log.Info(resp, "RESPONSE")
	session.Send(code, resp)
}
