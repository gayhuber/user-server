// Package session 新建会话信息
package session

import (
	"io"
	"net"
	"strings"
	"time"
)

// Accept 接收参数
func Accept(listener net.Listener) (net.Conn, error) {
	var tempDelay time.Duration

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}

		return conn, nil
	}
}

// Listen 监听端口
func Listen(network, address string, handler FuncHandler) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	return NewServer(listener).Run(handler)
}
