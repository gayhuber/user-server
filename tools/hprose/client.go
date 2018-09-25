package hprose

import (
	"errors"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
	"log"
	"sync"
	"time"
	"user-server/config"
)

var (
	// ADDRESS 服务地址
	ADDRESS = "tcp://127.0.0.1:2333"
	// TIMEOUT 超时设定
	TIMEOUT time.Duration = 10 * time.Second

	// ServiceMp 可用的函数组
	ServiceMp *ServiceMap
)

// ServiceMap 可用服务列表
type ServiceMap struct {
	handlers map[string]BaseClient
	RWLock   sync.RWMutex
	BaseService
}

// BaseService 远程调用服务所包含的方法
type BaseService struct {
	Callback func(module, class, method string, args map[string]interface{}) (interface{}, error) `name:"baseServer_callBack"`
}

// BaseClient 本地调用的集合
type BaseClient struct {
	Module string
	Class  string
	Func   string
}

func init() {
	host, ok := config.GetParam("rpc_base_host")
	if ok {
		ADDRESS = host
	}

	ServiceMp = &ServiceMap{}
	ServiceMp.handlers = make(map[string]BaseClient)

	log.Println("get method:", ServiceMp.handlers)

}

// AddMethod 添加方法
func (serv *ServiceMap) AddMethod(key string, body BaseClient) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("AddMethod error:", err)
		}
	}()

	if _, ok := serv.handlers[key]; !ok {
		serv.RWLock.Lock()
		serv.handlers[key] = body
		serv.RWLock.Unlock()
	}

	return nil
}

// GetClient 获取远程链接
func GetClient() rpc.Client {
	client := rpc.NewClient(ADDRESS)
	client.SetTimeout(TIMEOUT)
	return client
}

// GetBaseClient 获取 base 服务
func GetBaseClient() (*BaseService, rpc.Client) {
	client := GetClient()
	var base *BaseService
	client.UseService(&base)
	return base, client
}

// RemoteFunc 远程调用方法
func RemoteFunc(key string, args map[string]interface{}) (interface{}, error) {
	base, conn := GetBaseClient()
	defer conn.Close()

	serv, ok := ServiceMp.handlers[key]
	if !ok {
		msg := fmt.Sprintf("service: %s not found", key)
		return nil, errors.New(msg)
	}
	return base.Callback(serv.Module, serv.Class, serv.Func, args)
}
