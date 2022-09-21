package model

// NotifyModel 消息通知结构体
type NotifyModel struct {
	OrmModel
	CreatorModel

	TargetUserId int64  `json:"targetUserId" gorm:"not null;comment:推送目标用户id"`
	TypeCode     string `json:"typeCode" gorm:"not null;comment:消息类型"`
	Title        string `json:"title" gorm:"not null;comment:消息标题"`
	Content      string `json:"content" gorm:"comment:消息内容"`
	Uri          string `json:"uri" gorm:"comment:消息路由"`
}
