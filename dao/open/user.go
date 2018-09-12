package open

import (
	"time"
)

// DaoUser db_open 库下的 User 表
type DaoUser struct {
	ID           int    `gorm:"primary_key;type:bigint(20);not null;AUTO_INCREMENT"`
	OpenID       string `gorm:"type:varchar(32);not null;index:idx_open"`
	SyUID        int    `gorm:"column:sy_uid;type:bigint(20);not null;default 0"`
	Status       int    `gorm:"type:tinyint(4)"`
	Src          string
	Device       string
	Password     string
	PasswordSalt string `gorm:"type:varchar(32);not null"`
	CreateAt     time.Time
	UpdateAt     time.Time
	Token        string `gorm:"index:idx_token"`
	// 关联其他表内容,  foreignkey 指其他表的关联字段, AssociationForeignKey 为本地关联字段
	Info DaoUserInfo `gorm:"foreignkey:ID;AssociationForeignKey:ID"`
}

const (
	USER_STATUS_ACTIVE  int = 20
	USER_STATUS_PREPARE int = 10
)

// TableName 指定了这个 struct 依赖的表名
func (u DaoUser) TableName() string {
	return "tb_u_user"
}

// SaveNewUser 新增一个用户
func (u *DaoUser) SaveNewUser() (err error) {
	err = db.Create(u).Error

	// u.Info = DaoUserInfo{
	// 	Avatar: "http://img2.soyoung.com/user/1.png"
	// }
	return
}