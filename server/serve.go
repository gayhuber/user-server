// Package server server 这里放置了入口文件
package server

import (
	"fmt"
	"log"
	"user-server/lib/link"
	"user-server/lib/link/codec"
)

type AddReq struct {
	A, B int
}

type AddRsp struct {
	C int
}

// RequestParams 请求参数格式
type RequestParams struct {
	Route  string `json:"route"`
	LogID  string `json:"log_id"`
	Params interface{}
}

// Response 返回格式
type Response struct {
	Code int `json:"code"`
	Body map[string]interface{}
}

// JSON 返回的 JSON 数据
type JSON map[string]interface{}

// Hello 测试函数
func Hello() {
	addr := fmt.Sprintf("%s:%d", Conf.Main.Host, Conf.Main.Port)
	fmt.Println("this is server file", addr)
}

// Run 启动服务
func Run() {
	json := codec.Json()
	json.Register(RequestParams{})
	json.Register(Response{})
	json.Register(AddReq{})
	json.Register(AddRsp{})

	addr := fmt.Sprintf("%s:%d", Conf.Main.Host, Conf.Main.Port)
	server, err := link.Listen("tcp", addr, json, 0 /* sync send */, link.HandlerFunc(funcTransfer))
	checkErr(err)

	log.Println("server address is:", addr)
	server.Serve()
}

// 根据传入的 route 参数来使用不同的函数
func funcTransfer(session *link.Session) {
	for {
		req, err := session.Receive()
		checkErr(err)

		log.Println("server: get reuqest =>", req)
		route, _ := req.(map[string]interface{})
		log.Println("current route is:", route["route"])

		err = session.Send(&Response{
			Code: 200,
		})
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
