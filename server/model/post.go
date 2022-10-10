package model

// Post 帖子结构体
type Post struct {
	OrmBase
	Creator

	Title    string `json:"title" gorm:"varchar(80);not null;comment:帖子标题"`
	Contents []any  `json:"contents" gorm:"type:json;serializer:json;not null;comment:帖子内容"`

	TagCodes         []string        `json:"tagCodes" gorm:"type:json;serializer:json;comment:标签值集合"`
	Tags             []SimpleDict    `json:"tags,omitempty" gorm:"-"`
	ActivityRecordId *string         `json:"activityRecordId" gorm:"comment:活动id"`
	ActivityRecord   *ActivityRecord `json:"activityRecord,omitempty" gorm:"foreignKey:ActivityRecordId"`
	RecipeId         *string         `json:"recipeId" gorm:"comment:食谱id"`
	Recipe           *Recipe         `json:"recipe,omitempty" gorm:"foreignKey:RecipeId"`

	// 点赞/收藏过的用户
	LikeUsers    []User `json:"-" gorm:"many2many:post_like_users"`
	Liked        bool   `json:"liked" gorm:"-"`
	LikeCount    int64  `json:"likeCount" gorm:"-"`
	CollectUsers []User `json:"-" gorm:"many2many:post_collect_users"`
	Collected    bool   `json:"collected" gorm:"-"`
	CollectCount int64  `json:"collectCount" gorm:"-"`
}
