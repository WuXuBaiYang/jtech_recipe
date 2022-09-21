package model

// MenuMode 菜单信息结构
type MenuMode struct {
	OrmModel
	CreatorModel

	Contents []any `json:"contents" gorm:"type:json;serializer:json;not null;comment:菜单内容集合"`
	// 使用的是活动记录表中的id
	ActivityId *int64         `json:"activityId" gorm:"comment:活动id"`
	Activity   *ActivityModel `json:"activity" gorm:"-"`
}
