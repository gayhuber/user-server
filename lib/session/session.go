package session

import (
	"net"
	"sync"
	"sync/atomic"
)

var globalSessionID uint64

// Session 会话信息
type Session struct {
	conn     net.Conn
	request  interface{}
	response interface{}
	rwMutex  sync.RWMutex
	mutex    sync.Mutex
	id       uint64
}

// NewSession 新建会话
func NewSession(conn net.Conn) *Session {
	return &Session{
		conn: conn,
		id:   atomic.AddUint64(&globalSessionID, 1),
	}
}
