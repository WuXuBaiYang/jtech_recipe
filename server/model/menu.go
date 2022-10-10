package model

// Menu 菜单信息结构
type Menu struct {
	OrmBase
	Creator

	Contents   []any   `json:"contents" gorm:"type:json;serializer:json;not null;comment:菜单内容集合"`
	OriginId   *string `json:"originId" gorm:"comment:菜单复制来源id"`
	OriginMenu *Menu   `json:"originMenu,omitempty" gorm:"foreignKey:OriginId"`

	ActivityRecordId *string         `json:"activityRecordId" gorm:"comment:活动id"`
	ActivityRecord   *ActivityRecord `json:"activityRecord,omitempty" gorm:"foreignKey:ActivityRecordId"`
}
