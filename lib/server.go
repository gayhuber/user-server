package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"sync"
)

const (
	// Port 服务端接受的端口
	Port = ":6666"
)

/*
Open 返回一个有超时的TCP链接缓冲readwrite
net.Conn 实现了io.Reader  io.Writer  io.Closer接口
*/
func Open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	// Note that the local port is chosen on the fly. If the local port
	// must be a specific one, use DialTCP() instead.
	fmt.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// HandleFunc 请求处理函数类型
type HandleFunc func(*Session)

// TCPServer 服务端结构
type TCPServer struct {
	listener net.Listener
	// handlefunc是一个处理传入命令的函数类型。 它接收打包在一个读写器界面中的开放连接。
	handler map[string]HandleFunc

	// map不是线程安全的，所以需要读写锁控制
	m sync.RWMutex
}

// Session 会话信息
type Session struct {
	Conn     net.Conn
	Request  Request
	Response interface{}
	RwMutex  sync.RWMutex
	Mutex    sync.Mutex
	id       uint64
}

// NewTCPServer 初始化服务
func NewTCPServer() *TCPServer {
	return &TCPServer{
		handler: map[string]HandleFunc{},
	}
}

// AddHandleFunc 添加数据类型处理方法
func (serv *TCPServer) AddHandleFunc(name string, f HandleFunc) {
	serv.m.Lock()
	serv.handler[name] = f
	serv.m.Unlock()
}

// Request 的结构
type Request struct {
	Route  string `json:"route"`
	LogID  string `json:"log_id"`
	Params Params
}

// Params 参数的类型
type Params map[string]interface{}

// handleMessage 验证请求数据路由，并发送到对应处理函数
func (serv *TCPServer) handleMessage(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn),
		bufio.NewWriter(conn))
	defer conn.Close()

	by, err := rw.ReadBytes('\n')
	switch {
	case err == io.EOF:
		fmt.Println("读取完成.")
		return
	case err != nil:
		fmt.Println("读取出错")
		return
	}

	log.Println("获取请求:", string(by))

	var params Params
	err = json.Unmarshal(by, &params)
	if err != nil {
		log.Println("json 格式不正确, error:", err)
	}
	fmt.Println("收到 map 类参数: ", params)

	session := Session{
		Conn: conn,
		Request: Request{
			Route: params["route"].(string),
			LogID: params["log_id"].(string),
		},
	}

	serv.m.RLock()
	defer serv.m.RUnlock()
	handleCmd, ok := serv.handler[session.Request.Route]

	if !ok {
		fmt.Println("找不到对应路由规则: ", session.Request.Route)
		return
	}

	//具体处理链接数据
	handleCmd(&session)
}

// Listen 监听端口
func (serv *TCPServer) Listen() error {
	var err error
	serv.listener, err = net.Listen("tcp", Port)
	if err != nil {
		return errors.Wrap(err, "TCP服务无法监听在端口"+Port)
	}
	fmt.Println(" 服务监听成功：", serv.listener.Addr().String())
	for {
		conn, err := serv.listener.Accept()
		if err != nil {
			fmt.Println("新请求监听失败!")
			continue
		}
		// 开始处理新链接数据
		go serv.handleMessage(conn)
	}

}
