package models

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
func NewSession(token, Type string) *SessionModel {
	sm := &SessionModel{
		Token: token,
		Type:  Type,
	}
	return sm
}

func (sess *SessionModel) start() {

}

func (sess *SessionModel) store(info map[string]interface{}) {

}

// StoreSession 存储 session信息
func StoreSession(info map[string]interface{}) {

}
