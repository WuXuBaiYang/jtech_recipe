package model

// Comment 评论结构体
type Comment struct {
	OrmBase
	Creator

	PId      string `json:"pId" gorm:"comment:父级id"`
	TypeCode string `json:"typeCode" gorm:"comment:评论类型"`
	Content  string `json:"content" gorm:"varchar(300);not null;comment:评论文本内容"`

	// 点赞过的用户
	LikeUsers []User `json:"-" gorm:"many2many:comment_like_users"`
	Liked     bool   `json:"liked" gorm:"-"`
	LikeCount int64  `json:"likeCount" gorm:"-"`
}

// CommentType 评论类型枚举
type CommentType string

const (
	PostComment     CommentType = "11"
	RecipeComment   CommentType = "12"
	MenuComment     CommentType = "13"
	ActivityComment CommentType = "14"
)
