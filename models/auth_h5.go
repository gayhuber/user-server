package models

import (
	"fmt"
	"user-server/dao/open"
	"user-server/lib"
	"user-server/tools"
)

type H5Auth struct {
	Src      string
	Ext      map[string]interface{}
	uniqueID string
}

func (auth *H5Auth) getName() string {
	return "H5Auth"
}

// 注册一个新的用户
func (auth *H5Auth) register() (code int, obj interface{}) {
	openID, salt := tools.GenerageUniqueID(auth.Src)
	user := open.DaoUser{
		OpenID:       openID,
		Src:          auth.Src,
		Status:       open.USER_STATUS_PREPARE,
		PasswordSalt: salt,
		Token:        generateToken(openID, auth.Src, salt),
	}

	err := user.SaveNewUser()

	if err != nil {
		return 400, lib.H{
			"result": "no new user",
			"src":    auth.Src,
			"other":  auth.Ext,
			"msg":    err,
		}
	}

	fmt.Println(user)

	return 200, lib.H{
		"result": "this is h5 handler",
		"src":    auth.Src,
		"other":  auth.Ext,
	}
}

func (auth *H5Auth) login() (code int, obj interface{}) {
	return
}

func (auth *H5Auth) info() (code int, obj interface{}) {
	return
}

func (auth *H5Auth) home() (code int, obj interface{}) {
	return
}

func (auth *H5Auth) setParams(params map[string]interface{}) {
	auth.Src = params["src"].(string)
	if auth.Src != "fanli" {
		auth.Ext = params
	}

	return
}

func generateToken(openID, src, salt string) string {
	key := []byte(openID + src + salt)
	return tools.MD5String(key)
}
