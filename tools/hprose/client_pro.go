package hprose

import (
	"errors"
	"fmt"
	"github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

// BaseServiceV2 远程调用服务所包含的方法, 传参改成数组模式
type BaseServiceV2 struct {
	Callback func(module, class, method string, args []interface{}) (interface{}, error) `name:"baseServer_callBack"`
}

// 根据serverName 获取ip端口
func getServerHost(serverName string) (ipPort string, err error) {
	servers := []string{"register_center1:2181", "register_center2:2181", "register_center3:2181", "register_center4:2181", "register_center5:2181"}
	c, _, err := zk.Connect(servers, time.Second*10)
	defer c.Close()

	if err != nil {
		return
	}

	path := fmt.Sprintf("/server_list/%s", serverName)

	children, _, _, err := c.ChildrenW(path)

	if err != nil {
		return
	}

	if len(children) > 0 {
		ipPort = children[0]
		return
	}
	err = errors.New("don't get ipPort")
	return

}

// GetClientByServerName 通过服务名获取rpc客户端
func GetClientByServerName(serverName string) (*BaseServiceV2, rpc.Client) {
	ipPort, err := getServerHost(serverName)

	if err != nil {
		panic(err)
	}
	address := fmt.Sprintf("tcp://%s", ipPort)
	client := rpc.NewClient(address)
	client.SetTimeout(7 * time.Second)

	var base *BaseServiceV2
	client.UseService(&base)
	return base, client
}

// RemoteFuncPro 升级版的根据zk获取服务地址, 并调用方法
func RemoteFuncPro(serverName, moduleName, className, funcName string, args []interface{}) (interface{}, error) {
	start := time.Now()
	base, conn := GetClientByServerName(serverName)
	defer conn.Close()

	zkTime := time.Since(start)

	result, err := base.Callback(moduleName, className, funcName, args)

	log.Println("Callback result: ", result, err)

	if err != nil {
		return nil, err
	}

	baseTime := time.Since(start)
	// 将map[interface{}]interface{} 数据类型转换成可被json解析的格式
	w := io.NewWriter(true)
	w.Serialize(result)
	reader := io.NewReader(w.Bytes(), true)
	reader.JSONCompatible = true
	var p interface{}
	reader.Unserialize(&p)

	decodeTime := time.Since(start)

	log.Println("RemoteFuncPro execute time: ", zkTime, baseTime, decodeTime)

	return p, err
}
