package main

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
	"user-server/lib"
)

func server() error {
	TCPServer := lib.NewTCPServer()
	TCPServer.AddHandleFunc("test/index", HandleJSON)
	TCPServer.AddHandleFunc("activity/contact/index", HandleJSON)

	// 开始监听
	return TCPServer.Listen()
}

func main() {
	err := server()
	if err != nil {
		fmt.Println("Error:", errors.WithStack(err))
	}
	time.Sleep(time.Second * 100)
}

// HandleJSON 处理 json 文件
func HandleJSON(session *lib.Session) {

	fmt.Println("hello", session.Request)

	session.Log.Info("tesetse")

	session.Send(200, lib.H{
		"message": "admin hi",
		"param":   "example",
		"raw":     session.Request.Params,
	})
}
