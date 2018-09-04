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
	by, err := rw.ReadBytes('\n')

	if err != nil {
		return nil, err
	}

	log.Println("获取请求:", string(by))

	var params Params
	err = json.Unmarshal(by, &params)
	if err != nil {
		log.Println("json 格式不正确, error:", err)
	}
	fmt.Println("收到 map 类参数: ", params)

	session := &Session{
		Conn: conn,
		Request: Request{
			Route:  params["route"].(string),
			LogID:  params["log_id"].(string),
			Params: params,
		},
		id: atomic.AddUint64(&globalSessionID, 1),
	}

	session.Log = logs.NewLog(session.Request.LogID, Conf.Log.Path, Conf.Log.Mode)

	return session, nil
}

// Send 目前是回复 json
func (s *Session) Send(code int, obj interface{}) {
	var resp Response
	resp.Code = code
	resp.Data = obj

	jsons, err := json.Marshal(resp)
	if err != nil {
		log.Println("转换出错:", err)
	}

	s.Conn.Write(jsons)
	// 输入结束标志
	s.Conn.Write([]byte(Conf.Main.EOF))

	log.Println("输出 json:", string(jsons))
}
