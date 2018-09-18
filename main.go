package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"time"
	"user-server/lib"
	"user-server/models"
)

func server() error {
	TCPServer := lib.NewTCPServer()
	TCPServer.AddHandleFunc("demo/test", HandleJSON)
	TCPServer.AddHandleFunc("activity/contact/index", HandleJSON)
	TCPServer.AddHandleFunc("auth/user/register", models.AuthRegister)
	TCPServer.AddHandleFunc("auth/user/login", models.AuthLogin)

	// 开始监听
	return TCPServer.Listen()
}

func main() {
	err := server()
	if err != nil {
		fmt.Println("Error:", errors.WithStack(err))
		os.Exit(0)
	}
	time.Sleep(time.Second * 100)
}

// HandleJSON 处理 json 文件
func HandleJSON(session *lib.Session) {

	fmt.Println("hello", session.Request)

	session.Log.Info("tesetse")

	session.Send(200, lib.H{
		"message": "this is from server",
		"param":   "example",
		"raw":     session.Request.Params,
	})
}
