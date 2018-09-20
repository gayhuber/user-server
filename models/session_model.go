package models

import (
	"errors"
	"user-server/dao/open"
	"user-server/tools"
	logs "user-server/tools/loghandler"
)

const (
	SUBMIT_LIFE    int    = 864000
	SESSION_PREFIX string = "session_"
)

// SessionModel 维护用户登录信息的模型
type SessionModel struct {
	SessionID   string
	SessionInfo map[string]interface{}
	Type        string
	Token       string
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

// strore 将用户信息存入缓存
func (sess *SessionModel) store(info map[string]interface{}) {
	sess.SessionInfo = info
	StoreSession(sess.getKey(), info)
}

// getKey 组件缓存 key
func (sess *SessionModel) getKey() string {
	return SESSION_PREFIX + sess.Token
}

// StoreSession 存储 session 信息
func StoreSession(key string, info map[string]interface{}) {
	redis := tools.GetRedis()
	defer redis.Close()
	infoByte, _ := tools.JSONEncode(info)
	err := redis.Set(key, infoByte, SUBMIT_LIFE).Error
	if err != nil {
		logs.Error(err, "SESSION_MODEL_SAVE_ERROR")
	}
}
