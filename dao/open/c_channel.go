package open

import (
	"time"
)

// CChannel 用户渠道信息表
type CChannel struct {
	ID          int `gorm:"primary_key;type:int(11);not null;AUTO_INCREMENT"`
	Type        string
	FstKey      string
	SecKey      string
	Status      int
	QueryURL    string `gorm:"column:query_url"`
	SubmitURL   string `gorm:"column:submit_url"`
	CreateTime  time.Time
	UpdateTime  time.Time
	XyPlatID    int `gorm:"column:xy_plat_id"`
	Ext         string
	ChannelName string
}

// TableName 指定了这个 struct 依赖的表名
func (c CChannel) TableName() string {
	return "tb_u_info"
}
