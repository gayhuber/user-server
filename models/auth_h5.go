package models

import (
	"time"
	"user-server/dao/open"
	"user-server/lib"
	"user-server/tools"
)

type H5Auth struct {
	Src      string
	Ext      map[string]interface{}
	uniqueID string
	OpenID   string
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
		CreateAt:     time.Now(),
	}

	err := user.SaveNewUser()

	if err != nil {
		return 500, lib.H{
			"msg": "数据库插入失败",
		}
	}

	syUID, isRegister := tools.RegisteSoyoungHalf(user.OpenID, user.Src, user.Info.Avatar, user.Info.Nickname)
	if !isRegister {
		return 500, lib.H{
			"msg": "注册新氧半账号失败",
		}
	}

	err = user.FlashSyUID(syUID)
	if err != nil {
		return 500, lib.H{
			"msg": "数据库更新失败",
		}
	}

	// 正常返回的内容
	return 200, lib.H{
		"open_id":       user.OpenID,
		"src":           user.Src,
		"status":        user.Status,
		"password_salt": user.PasswordSalt,
		"token":         user.Token,
		"id":            user.ID,
		"sy_uid":        user.SyUID,
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
	if src, ok := params["src"]; ok {
		auth.Src = src.(string)
	}
	if openID, ok := params["open_id"]; ok {
		auth.OpenID = openID.(string)
	}
	auth.Ext = params

	return
}

func generateToken(openID, src, salt string) string {
	key := []byte(openID + src + salt)
	return tools.MD5String(key)
}
