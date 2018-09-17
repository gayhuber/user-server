package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	logs "user-server/tools/loghandler"
	// 引入配置文件
	"user-server/config"
)

var globalSessionID uint64

// Session 会话信息
type Session struct {
	Conn     net.Conn
	Request  Request
	Response Response
	RwMutex  sync.RWMutex
	Mutex    sync.Mutex
	id       uint64
	Log      logs.Logger
}

// NewSession 新建 session
func NewSession(rw *bufio.ReadWriter, conn net.Conn) (*Session, error) {
	// 错误集中处理
	defer paincHandler(conn)

	by, err := rw.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var params Params
	err = json.Unmarshal(by, &params)
	if err != nil {
		logs.Error("json 格式不正确, error: %s", err)
		panic(err.Error())
	}

	session := &Session{
		Conn: conn,
		Request: Request{
			Route:  params["route"].(string),
			LogID:  params["log_id"].(string),
			Params: params,
		},
		id: atomic.AddUint64(&globalSessionID, 1),
	}

	session.Log = logs.NewLog(session.Request.LogID)
	session.Log.Info(params, "GET_PARAMS")

	return session, nil
}

// Send 目前是回复 json
func (s *Session) Send(code int, obj interface{}) {
	resp := ResponseHandler(code, obj)

	jsons, err := json.Marshal(resp)
	if err != nil {
		log.Println("转换出错:", err)
	}

	s.Conn.Write(jsons)
	// 输入结束标志
	s.Conn.Write([]byte(config.Conf.Main.EOF))

	s.Log.Info(string(jsons), "RESPONSE")

}

func paincHandler(conn net.Conn) {
	session := &Session{
		Conn: conn,
	}
	if err := recover(); err != nil {
		logs.Error(err)
		str := fmt.Sprint(err)
		session.Send(500, str)
	}
}
