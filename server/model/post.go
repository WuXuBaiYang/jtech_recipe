package model

// PostModel 帖子结构体
type PostModel struct {
	OrmModel
	CreatorModel

	Title    string `json:"title" gorm:"varchar(80);not null;comment:帖子标题"`
	Contents []any  `json:"contents" gorm:"type:json;serializer:json;not null;comment:帖子内容"`

	TagCodes []string        `json:"tagCodes" gorm:"type:json;serializer:json;comment:标签值集合"`
	Tags     []RespDictModel `json:"tags" gorm:"-"`

	// 使用的是活动记录表中的id
	ActivityId *int64         `json:"activityId" gorm:"comment:活动id"`
	Activity   *ActivityModel `json:"activity" gorm:"-"`
	RecipeId   *int64         `json:"recipeId" gorm:"comment:食谱id"`
	Recipe     *RecipeModel   `json:"recipe" gorm:"-"`

	// 浏览/点赞/收藏过的用户
	ViewUsers    []UserModel `json:"-" gorm:"many2many:post_view_users"`
	LikeUsers    []UserModel `json:"-" gorm:"many2many:post_like_users"`
	CollectUsers []UserModel `json:"-" gorm:"many2many:post_collect_users"`
}

// PostCommentModel 帖子评论结构体
type PostCommentModel struct {
	OrmModel
	CreatorModel

	PId     int64  `json:"pId" gorm:"comment:父级id"`
	Content string `json:"content" gorm:"varchar(300);not null;comment:帖子评论文本内容"`

	// 点赞过的用户
	LikeUsers []UserModel `json:"-" gorm:"many2many:post_comment_like_users"`
}

// PostCommentReplayModel 帖子评论回复结构体
type PostCommentReplayModel struct {
	OrmModel
	CreatorModel

	PId     int64  `json:"pId" gorm:"comment:父级id"`
	Content string `json:"content" gorm:"varchar(300);not null;comment:帖子评论回复文本内容" `

	// 点赞过的用户
	LikeUsers []UserModel `json:"-" gorm:"many2many:post_comment_replay_like_users"`
}
