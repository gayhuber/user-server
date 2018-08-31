// Package session 新建会话信息
package session

import (
	"net"
)

// Server 一个服务主体
type Server struct {
	listener net.Listener
}

// FuncHandler 注入方法
type FuncHandler func(session *Session)

// NewServer 创建一个新的服务
func NewServer(l net.Listener) *Server {
	return &Server{
		listener: l,
	}
}

// Listener 获取监听
func (server *Server) Listener() net.Listener {
	return server.listener
}

// Run 启动服务
func (server *Server) Run(handler FuncHandler) error {
	for {
		conn, err := Accept(server.listener)
		if err != nil {
			conn.Close()
			return err
		}
		// 每来一个请求都将产生一个 session
		session := NewSession(conn)
		go handler(session)
	}
}
