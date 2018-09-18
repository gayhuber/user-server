package open

import ()

// DaoUserInfo 用户信息表
type DaoUserInfo struct {
	ID        int    `gorm:"primary_key;type:bigint(20);not null;AUTO_INCREMENT"`
	Avatar    string `json:"avatar"`
	Mobile    string
	RealName  string
	Gender    int
	IDCode    string `gorm:"column:id_code"`
	Status    int
	Nickname  string `json:"nickname"`
	ChannelID string `gorm:"column:channel_uid"`
	Ext       string
}

// TableName 指定了这个 struct 依赖的表名
func (i DaoUserInfo) TableName() string {
	return "tb_u_info"
}
