// Package server config 这里放置了相关的所有配置
package server

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 整个服务的配置结构
type Config struct {
	Main main
}

type main struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	EOF  string `json:"eof"`
}

var (
	// Conf 运行时的配置
	Conf *Config
)

func init() {
	InitConfig()
	fmt.Println("hello")
	fmt.Println("host address:", viper.Get("main.host"))
	fmt.Println(Conf)
}

// NewConfig 生成一个新的配置文件
func NewConfig() *Config {
	return &Config{}
}

// InitConfig 初始化配置
func InitConfig() (err error) {
	Conf = NewConfig()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file : %s ", err))
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unable to decode into structL %s ", err))
	}
	return nil
}
