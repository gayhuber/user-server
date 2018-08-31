package server

import (
	"fmt"
	"log"
	"user-server/lib/link"
	"user-server/lib/link/codec"
)

type req struct {
	Route string `json:"route"`
	LogID string `json:"log_id"`
	IP    string `json:"ip"`
	Num   int    `json:"num"`
	Token string `json:"token"`
}

// Test 测试文件
func Test() {
	json := codec.Json()
	json.Register(req{})

	addr := fmt.Sprintf("%s:%d", Conf.Main.Host, Conf.Main.Port)
	client, err := link.Dial("tcp", addr, json, 0)
	checkErr(err)
	clientSessionLoop(client)

}

func clientSessionLoop(session *link.Session) {
	for i := 0; i < 3; i++ {
		err := session.Send(&req{
			Route: "/user/log/bin",
			LogID: "aabbcc",
			IP:    "127.0.0.1",
			Num:   i,
			Token: "abcdefg",
		})
		checkErr(err)
		log.Printf("Send: %d", i)

		rsp, err := session.Receive()
		checkErr(err)
		log.Printf("Receive: %d", rsp)
	}
}
