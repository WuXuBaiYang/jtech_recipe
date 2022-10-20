package model

// Menu 菜单信息结构
type Menu struct {
	OrmBase
	Creator

	Title            string          `json:"title" gorm:"comment:菜单标题"`
	Contents         []any           `json:"contents" gorm:"type:json;serializer:json;not null;comment:菜单内容集合"`
	OriginId         *string         `json:"originId" gorm:"comment:菜单复制来源id"`
	OriginMenu       *Menu           `json:"originMenu,omitempty" gorm:"foreignKey:OriginId"`
	TagCodes         []string        `json:"tagCodes" gorm:"type:json;serializer:json;comment:标签值集合"`
	Tags             []SimpleDict    `json:"tags" gorm:"-"`
	ActivityRecordId *string         `json:"activityRecordId" gorm:"comment:活动id"`
	ActivityRecord   *ActivityRecord `json:"activityRecord,omitempty" gorm:"foreignKey:ActivityRecordId"`

	// 点赞/收藏过的用户
	LikeUsers    []User `json:"-" gorm:"many2many:menu_like_users"`
	Liked        bool   `json:"liked" gorm:"-"`
	LikeCount    int64  `json:"likeCount" gorm:"-"`
	CollectUsers []User `json:"-" gorm:"many2many:menu_collect_users"`
	Collected    bool   `json:"collected" gorm:"-"`
	CollectCount int64  `json:"collectCount" gorm:"-"`
}
