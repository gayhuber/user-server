package models

import (
	"encoding/json"
	"errors"
	"user-server/dao/open"
	"user-server/tools"
	logs "user-server/tools/loghandler"
)

const (
	// SUBMITLIFE 缓存周期
	SUBMITLIFE int = 864000
	// SESSIONPREFIX 缓存前缀
	SESSIONPREFIX string = "session_"
)

// SessionModel 维护用户登录信息的模型
type SessionModel struct {
	SessionID   string
	SessionInfo map[string]interface{}
	Type        string
	Token       string
}

type tokenBody struct {
	UID         string `json:"uid"`
	Token       string `json:"token"`
	LoginMobile string
	LoginName   string
	NewUser     int    `json:"new_user"`
	XYToken     string `json:"xy_token"`
	Avatar      string
	Gendre      string
}

// NewSession 获取 session
//@token 用户登录 token
//@Type 登录类型, 目前可选有 h5, mobile
func NewSession(token, Type string) *SessionModel {
	sm := &SessionModel{
		Token: token,
		Type:  Type,
	}
	return sm
}

// info 这里虽然是获取 session 信息, 但是如果缓存中没有那么也要现存一份进来
func (sess *SessionModel) info() (res []byte, err error) {
	if len(sess.SessionInfo) > 0 {
		// 这里是为了把 map[string]interface{} 转成 []byte 用的一个迂回
		by, _ := tools.JSONEncode(sess.SessionInfo)
		return by, nil
	}

	redis := tools.GetRedis()
	defer redis.Close()

	if redis.Exist(sess.getKey()) {
		resp, err := redis.Get(sess.getKey())
		return resp.([]byte), err
	}

	// redis 中找不到的从数据库找, 再存一份到 redis
	if sess.Type != "mobile" {
		user := open.DaoUser{}
		info, err := user.FindByToken(sess.Token)
		if err != nil {
			return []byte{}, err
		}

		res, err := tools.JSONEncode(info)
		sess.store(info)
		return res, err
	}

	return []byte{}, errors.New("nothing was happend in session")
}

func (sess *SessionModel) infoStruct() (tb tokenBody, err error) {
	resp, err := sess.info()
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &tb)
	return
}

// strore 将用户信息存入缓存
func (sess *SessionModel) store(info map[string]interface{}) {
	sess.SessionInfo = info
	if sess.Type == "mobile" {
		StoreSession(sess.getKey(), info, tools.CalculateTTL())
	}
	StoreSession(sess.getKey(), info)
}

// getKey 组件缓存 key
func (sess *SessionModel) getKey() string {
	return SESSIONPREFIX + sess.Token
}

// StoreSession 存储 session 信息
func StoreSession(key string, info map[string]interface{}, ttl ...int) {
	redis := tools.GetRedis()
	defer redis.Close()
	infoByte, _ := tools.JSONEncode(info)

	var exp int
	if len(ttl) > 0 {
		exp = ttl[0]
	} else {
		exp = SUBMITLIFE
	}

	err := redis.Set(key, infoByte, exp).Error
	if err != nil {
		logs.Error(err, "SESSION_MODEL_SAVE_ERROR")
	}
}
