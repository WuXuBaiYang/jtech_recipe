package model

// Notification 消息通知结构体
type Notification struct {
	OrmModel
	CreatorModel

	TargetUserId uint   `json:"targetUserId" gorm:"not null;comment:推送目标用户id"`
	Type         int64  `json:"type" gorm:"not null;comment:消息类型（0系统消息|1关注对象发帖）"`
	Title        string `json:"title" gorm:"not null;comment:消息标题"`
	Content      string `json:"content" gorm:"comment:消息内容"`
	Uri          string `json:"uri" gorm:"comment:消息路由"`
}
