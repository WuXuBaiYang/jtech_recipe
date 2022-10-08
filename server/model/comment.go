package model

// Comment 评论结构体
type Comment struct {
	OrmBase
	Creator

	PId     string `json:"pId" gorm:"comment:父级id"`
	Content string `json:"content" gorm:"varchar(300);not null;comment:评论文本内容"`

	// 点赞过的用户
	LikeUsers []User `json:"-" gorm:"many2many:comment_like_users"`
	Liked     bool   `json:"liked" gorm:"-"`
	LikeCount int64  `json:"likeCount" gorm:"-"`
}
