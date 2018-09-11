// Package dao  这里写 Dao 文件
package dao

import (
	"github.com/jinzhu/gorm"
	// 引入 gorm 的 mysql 支持
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"log"
	"os"
	// 引入配置文件
	"fmt"
	"user-server/config"
)

// GetDB 获取一个链接信息
func GetDB(dbName string) *gorm.DB {
	conf, err := GetConf(dbName)
	fmt.Println(conf)
	if err != nil {
		log.Println(err)
	}

	// confStrDemo := "root:123123@tcp(127.0.0.1:33060)/db_open?charset=utf8&parseTime=True&loc=Local"
	confStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database)

	db, err := gorm.Open("mysql", confStr)
	if err != nil {
		log.Println("orm open error", err, confStr)
		os.Exit(0)
	}

	// 设置最大链接数
	db.DB().SetMaxOpenConns(conf.MaxOpenConn)
	// 设置最大闲置数
	db.DB().SetMaxIdleConns(conf.MaxIdleConn)
	return db
}

// GetConf 获取数据库配置
func GetConf(key string) (conf config.DBConfig, err error) {
	conf, ok := config.Conf.DB[key]
	if !ok {
		err = errors.New("db config not found")
		return
	}

	return
}
