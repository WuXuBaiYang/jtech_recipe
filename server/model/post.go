package model

// Post 帖子结构体
type Post struct {
	OrmBase
	Creator

	Title    string `json:"title" gorm:"varchar(80);not null;comment:帖子标题"`
	Contents []any  `json:"contents" gorm:"type:json;serializer:json;not null;comment:帖子内容"`

	TagCodes []string     `json:"tagCodes" gorm:"type:json;serializer:json;comment:标签值集合"`
	Tags     []SimpleDict `json:"tags" gorm:"-"`

	// 使用的是活动记录表中的id
	ActivityRecordId *int64          `json:"activityId" gorm:"not null;comment:活动id"`
	ActivityRecord   *ActivityRecord `json:"activity" gorm:"foreignKey:ActivityRecordId"`
	RecipeId         *int64          `json:"recipeId" gorm:"comment:食谱id"`
	Recipe           *Recipe         `json:"recipe" gorm:"foreignKey:RecipeId"`

	// 浏览/点赞/收藏过的用户
	ViewUsers    []User `json:"-" gorm:"many2many:post_view_users"`
	Viewed       bool   `json:"viewed" gorm:"-"`
	ViewCount    int64  `json:"viewCount" gorm:"-"`
	LikeUsers    []User `json:"-" gorm:"many2many:post_like_users"`
	Liked        bool   `json:"liked" gorm:"-"`
	LikeCount    int64  `json:"likeCount" gorm:"-"`
	CollectUsers []User `json:"-" gorm:"many2many:post_collect_users"`
	Collected    bool   `json:"collected" gorm:"-"`
	CollectCount int64  `json:"collectCount" gorm:"-"`
}

// PostComment 帖子评论结构体
type PostComment struct {
	OrmBase
	Creator

	PId     int64  `json:"pId" gorm:"comment:父级id"`
	Content string `json:"content" gorm:"varchar(300);not null;comment:帖子评论文本内容"`

	// 点赞过的用户
	LikeUsers []User `json:"-" gorm:"many2many:post_comment_like_users"`
	Liked     bool   `json:"liked" gorm:"-"`
	LikeCount int64  `json:"likeCount" gorm:"-"`
}

// PostCommentReplay 帖子评论回复结构体
type PostCommentReplay struct {
	OrmBase
	Creator

	PId     int64  `json:"pId" gorm:"comment:父级id"`
	Content string `json:"content" gorm:"varchar(300);not null;comment:帖子评论回复文本内容" `

	// 点赞过的用户
	LikeUsers []User `json:"-" gorm:"many2many:post_comment_replay_like_users"`
	Liked     bool   `json:"liked" gorm:"-"`
	LikeCount int64  `json:"likeCount" gorm:"-"`
}
