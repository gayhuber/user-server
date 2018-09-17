// Package config 这里放置了相关的所有配置
package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// Config 整个服务的配置结构
type Config struct {
	Main   mainConfig
	Log    logConfig
	DB     map[string]DBConfig `json:"db"`
	Params map[string]string   `json:"params"`
}

//DBConfig 数据库配置文件
type DBConfig struct {
	Database    string `json:"database"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	MaxOpenConn int    `json:"max_open_conn"`
	MaxIdleConn int    `json:"max_idle_conn"`
	Port        int    `json:"port"`
	Host        string `json:"host"`
}

type mainConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	EOF  string `json:"eof"`
}

type logConfig struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
}

var (
	// Conf 运行时的配置
	Conf *Config

	h bool
	c string
)

// 初始化加载配置文件
func init() {
	// 错误处理
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("启动失败 ", err)
			os.Exit(0)
		}
	}()

	parseParams()
	InitConfig()
}

// parsePrams 解析参数
func parseParams() {
	flag.BoolVar(&h, "h", false, "get help")
	flag.StringVar(&c, "c", ".", "set config path")
	flag.Parse()
	if h {
		usage()
		os.Exit(0)
	}
}

// 参数提示
func usage() {
	fmt.Println(`
Usage: app  [-c filePath]

Options:
`)
	flag.PrintDefaults()
}

// NewConfig 生成一个新的配置文件
func NewConfig() *Config {
	return &Config{}
}

// InitConfig 初始化配置
func InitConfig() (err error) {

	viper.SetConfigName("config")
	viper.AddConfigPath(c)
	viper.SetConfigType("json")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%s ", err))
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unable to decode into structL %s ", err))
	}
	return nil
}
