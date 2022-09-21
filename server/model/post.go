package model

// PostModel 帖子结构体
type PostModel struct {
	OrmModel
	CreatorModel

	Title    string `json:"title" gorm:"varchar(80);not null;comment:帖子标题"`
	Contents []any  `json:"contents" gorm:"type:json;serializer:json;not null;comment:帖子内容"`

	TagCodes []string         `json:"tagCodes" gorm:"type:json;serializer:json;comment:标签值集合"`
	Tags     *[]RespDictModel `json:"tags,omitempty" gorm:"-"`
}

// PostCommentModel 帖子评论结构体
type PostCommentModel struct {
	OrmModel
	CreatorModel

	PId     int64  `json:"pId" gorm:"comment:父级id"`
	Content string `json:"content" gorm:"varchar(300);not null;comment:帖子评论文本内容"`
}

// PostCommentReplayModel 帖子评论回复结构体
type PostCommentReplayModel struct {
	OrmModel
	CreatorModel

	PId     int64  `json:"pId" gorm:"comment:父级id"`
	Content string `json:"content" gorm:"varchar(300);not null;comment:帖子评论回复文本内容" `
}
