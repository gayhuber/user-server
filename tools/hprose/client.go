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
	TIMEOUT time.Duration = 7 * time.Second
	// ServiceMp 可用的函数组
	ServiceMp *ServiceMap
)

const (
	// BaseServerKey base 服务在 config 中的 key 值
	BaseServerKey = "rpc_base_host"
	// SmsServerkey sms 服务在 config 中的 key 值
	SmsServerkey = "rpc_sms_host"
)

// ServiceMap 可用服务列表
type ServiceMap struct {
	handlers map[string]BaseClient
	RWLock   sync.RWMutex
	BaseService
	serverMap map[string]string
}

// BaseService 远程调用服务所包含的方法
type BaseService struct {
	Callback func(module, class, method string, args map[string]interface{}) (interface{}, error) `name:"baseServer_callBack"`
}

// BaseClient 本地调用的集合
type BaseClient struct {
	Server string // 所属 server 名, eg: base-server, sms-server,etc...
	Module string
	Class  string
	Func   string
}

func init() {
	basehost, ok := config.GetParam(BaseServerKey)
	if ok {
		ADDRESS = basehost
	}
	smshost, ok := config.GetParam(SmsServerkey)
	if !ok {
		smshost = ADDRESS
	}

	ServiceMp = &ServiceMap{}
	ServiceMp.handlers = make(map[string]BaseClient)
	ServiceMp.serverMap = make(map[string]string)

	// 注册 base 服务地址
	ServiceMp.serverMap[BaseServerKey] = basehost
	// 注册 sms 服务地址
	ServiceMp.serverMap[SmsServerkey] = smshost
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

// GetClientByKey 通过方法名获取 client 链接
func (cli *BaseClient) GetClientByKey(key string) (base *BaseService, client rpc.Client) {
	if cli.Server == "" {
		client = GetClient(ADDRESS)
	}
	host, ok := ServiceMp.serverMap[cli.Server]
	if !ok {
		client = GetClient(ADDRESS)
	} else {
		client = GetClient(host)
	}
	client.UseService(&base)
	return base, client
}

// GetClient 获取远程链接
func GetClient(address string) rpc.Client {
	client := rpc.NewClient(address)
	client.SetTimeout(TIMEOUT)
	return client
}

// GetBaseClient 获取 base 服务
func GetBaseClient() (*BaseService, rpc.Client) {
	client := GetClient(ADDRESS)
	var base *BaseService
	client.UseService(&base)
	return base, client
}

// RemoteFunc 远程调用方法
func RemoteFunc(key string, args map[string]interface{}, try ...int) (interface{}, error) {
	base, conn := GetBaseClient()
	defer conn.Close()

	serv, ok := ServiceMp.handlers[key]
	if !ok {
		msg := fmt.Sprintf("service: %s not found", key)
		return nil, errors.New(msg)
	}
	result, err := base.Callback(serv.Module, serv.Class, serv.Func, args)

	_, ok = result.([]interface{})
	if ok {
		errMsg := fmt.Sprintf("result return an error format data: %+v", result)
		err = errors.New(errMsg)
	}

	if err != nil {
		tryTime := 0
		if len(try) > 0 {
			tryTime = try[0]
		}
		fmt.Printf("base.Callback exec err, error: %+v , try: %+v \n", err, tryTime)
		time.Sleep(time.Microsecond * 500)
		if tryTime < 3 {
			tryTime++
			return RemoteFunc(key, args, tryTime)
		}
		return result, err
	}
	return result, err
}

// RemoteFuncPlus 远程调用方法(新)
func RemoteFuncPlus(key string, args map[string]interface{}, try ...int) (interface{}, error) {

	serv, ok := ServiceMp.handlers[key]
	if !ok {
		msg := fmt.Sprintf("service: %s not found", key)
		return nil, errors.New(msg)
	}

	base, conn := serv.GetClientByKey(key)
	defer conn.Close()

	result, err := base.Callback(serv.Module, serv.Class, serv.Func, args)

	_, ok = result.([]interface{})
	if ok {
		errMsg := fmt.Sprintf("result return an error format data: %+v", result)
		err = errors.New(errMsg)
	}

	if err != nil {
		tryTime := 0
		if len(try) > 0 {
			tryTime = try[0]
		}
		fmt.Printf("base.Callback exec err, error: %+v , try: %+v \n", err, tryTime)
		time.Sleep(time.Microsecond * 500)
		if tryTime < 3 {
			tryTime++
			return RemoteFunc(key, args, tryTime)
		}
		return result, err
	}
	return result, err
}
