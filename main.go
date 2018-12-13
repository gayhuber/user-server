package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"time"
	"user-server/lib"
	"user-server/models"
	"user-server/tools/hprose"
)

func server() error {
	TCPServer := lib.NewTCPServer()
	TCPServer.AddHandleFunc("demo/test", HandleJSON)
	TCPServer.AddHandleFunc("auth/user/info", models.AuthInfo)
	TCPServer.AddHandleFunc("auth/user/register", models.AuthRegister)
	TCPServer.AddHandleFunc("auth/user/login", models.AuthLogin)
	TCPServer.AddHandleFunc("auth/user/sms", models.MobileSms)
	TCPServer.AddHandleFunc("auth/user/home", models.MobileHome)
	TCPServer.AddHandleFunc("auth/user/account", models.MobileAccount)

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

	t1 := time.Now()

	body := session.Request.Params

	session.Log.Error(body["module"])

	res, _ := hprose.RemoteFuncPro(body["serverName"].(string), body["module"].(string), body["class"].(string), body["func"].(string), body["params"].(map[string]interface{}))

	session.Log.Info("tesetse")

	end := time.Since(t1)

	fmt.Println("execute time:", end)

	session.Send(200, res)
}
