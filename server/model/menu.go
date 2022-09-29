package model

// RecipeMenu 食谱菜单信息结构
type RecipeMenu struct {
	OrmBase
	Creator

	Contents []any `json:"contents" gorm:"type:json;serializer:json;not null;comment:菜单内容集合"`

	ActivityRecordId *string         `json:"activityId" gorm:"not null;comment:活动id"`
	ActivityRecord   *ActivityRecord `json:"activity" gorm:"foreignKey:ActivityRecordId"`
}
