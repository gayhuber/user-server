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

func transfer(tp string, params map[string]interface{}) (hd authHandler, err error) {
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

	session.Send(code, resp)

	// session.Send(200, lib.H{
	// 	"message": "this is from server",
	// 	"param":   "auth register",
	// 	"log_id":  session.Log.ID,
	// 	"raw":     session.Request.Params,
	// })
}

// func init() {
// 	handlerMap["h5"] = &h5Auth{}
// }
