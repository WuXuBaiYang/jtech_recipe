package model

// Post 帖子结构体
type Post struct {
	OrmModel
	CreatorModel
	Title    string `json:"title" gorm:"varchar(80);not null;comment:帖子标题"`
	Contents []any  `json:"contents" gorm:"type:json;serializer:json;not null;comment:帖子内容"`

	ViewUsers []User `json:"-" gorm:"many2many:post_view_users"`
	ViewCount int64  `json:"viewCount" gorm:"-"`
	Viewed    bool   `json:"viewed" gorm:"-"`

	LikeUsers    []User       `json:"-" gorm:"many2many:post_like_users"`
	LikeUserList []SimpleUser `json:"likeUserList,omitempty" gorm:"-"`
	LikeCount    int64        `json:"likeCount" gorm:"-"`
	Liked        bool         `json:"liked" gorm:"-"`

	CollectUsers []User `json:"-" gorm:"many2many:post_collect_users"`
	CollectCount int64  `json:"collectCount" gorm:"-"`
	Collected    bool   `json:"collected" gorm:"-"`

	Tags []PostTag `json:"tags" gorm:"foreignkey:PostID"`
}

// PostComment 帖子评论结构体
type PostComment struct {
	OrmModel
	CreatorModel
	PostID uint `json:"postId" gorm:"comment:所属帖子id"`

	Content string `json:"content" gorm:"varchar(300);not null;comment:帖子评论文本内容"`

	Replays     []PostCommentReplay `json:"replays" gorm:"foreignkey:CommentID"`
	ReplayCount int64               `json:"replayCount" gorm:"-"`

	LikeUsers []User `json:"-" gorm:"many2many:post_comment_like_users"`
	LikeCount int64  `json:"likeCount" gorm:"-"`
	Liked     bool   `json:"liked" gorm:"-"`
}

// PostCommentReplay 帖子评论回复结构体
type PostCommentReplay struct {
	OrmModel
	CreatorModel

	CommentID uint `json:"commentId" gorm:"comment:所属帖子评论id"`

	Content string `json:"content" gorm:"varchar(300);not null;comment:帖子评论回复文本内容" `

	LikeUsers []User `json:"-" gorm:"many2many:post_comment_replay_like_users"`
	LikeCount int64  `json:"likeCount" gorm:"-"`
	Liked     bool   `json:"liked" gorm:"-"`
}

// PostTag 帖子标签结构体
type PostTag struct {
	OrmModel
	CreatorModel

	PostID uint `json:"postId" gorm:"所属帖子id"`

	Name string `json:"name" gorm:"varchar(40);not null;unique;comment:帖子标签名称"`
}
