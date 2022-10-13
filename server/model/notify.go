package model

// Notify 消息通知结构体
type Notify struct {
	OrmBase

	FromUserId string      `json:"fromUserId" gorm:"not null;comment:消息来源用户id"`
	FromUser   *SimpleUser `json:"fromUser" gorm:"foreignKey:FromUserId"`
	ToUserId   string      `json:"toUserId" gorm:"not null;comment:消息目标用户id"`
	TypeCode   string      `json:"typeCode" gorm:"not null;comment:消息类型"`
	Title      string      `json:"title" gorm:"not null;comment:消息标题"`
	Content    string      `json:"content" gorm:"comment:消息内容"`
	Uri        string      `json:"uri" gorm:"comment:消息路由"`
}
