package open

import (
	"github.com/jinzhu/gorm"
	"log"
	"user-server/dao"
)

var db *gorm.DB

func init() {
	db = dao.GetDB("db_open")
	log.Println("connected mysql db_open")
}

// GetDB 获取当前库的数据连接
func GetDB() *gorm.DB {
	return db
}
