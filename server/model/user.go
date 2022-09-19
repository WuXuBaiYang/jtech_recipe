package model

import (
	"time"
)

// User 用户结构体
type User struct {
	OrmModel
	UserIM
	UserName string       `json:"userName" gorm:"varchar(20);not null;unique;comment:用户名（登录名）"`
	Password string       `json:"-" gorm:"size:255;not null;comment:密码（登录密码）"`
	Profile  *UserProfile `json:"profile,omitempty" gorm:"foreignkey:CreatorID"`

	// 浏览/点赞/收藏
	ViewPosts    []Post `json:"viewPosts,omitempty" gorm:"many2many:post_view_users"`
	LikePosts    []Post `json:"likePosts,omitempty" gorm:"many2many:post_like_users"`
	CollectPosts []Post `json:"collectPosts,omitempty" gorm:"many2many:post_collect_users"`

	// 发帖列表
	Posts []Post `json:"posts,omitempty" gorm:"foreignkey:CreatorID"`

	// 帖子评论/回复
	PostComments       []PostComment       `json:"postComments,omitempty" gorm:"foreignkey:CreatorID"`
	PostCommentReplays []PostCommentReplay `json:"postCommentReplays,omitempty" gorm:"foreignkey:CreatorID"`

	// 帖子/回复的点赞
	LikePostComments       []PostComment       `json:"likePostComments,omitempty" gorm:"many2many:post_comment_like_users"`
	LikePostCommentReplays []PostCommentReplay `json:"likePostCommentReplays,omitempty" gorm:"many2many:post_comment_replay_like_users"`

	// 创建的标签
	PostTags []PostTag `json:"PostTags,omitempty" gorm:"foreignkey:CreatorID"`

	// 关注列表
	SubscribeUsers []User `json:"subscribeUsers,omitempty" gorm:"many2many:user_subscribe;comment:用户订阅列表"`
}

// SimpleUser 简略用户结构体
type SimpleUser struct {
	ID uint `json:"id" gorm:"primarykey"`

	Profile *SimpleUserProfile `json:"profile,omitempty" gorm:"foreignkey:CreatorID;"`
}

// TableName 设置简略用户结构体表名
func (SimpleUser) TableName() string {
	return "user"
}

// UserProfile 用户信息结构体
type UserProfile struct {
	OrmModel
	CreatorID  uint      `json:"creatorId"`
	NickName   string    `json:"nickName" gorm:"varchar(20);comment:昵称"`
	Avatar     string    `json:"avatar" gorm:"varchar(40);comment:头像（只存储oss的key/id）"`
	Telephone  string    `json:"telephone" gorm:"varchar(20);default:null;unique;comment:手机号（只可以为空或者校验过的手机号格式）"`
	Bio        string    `json:"bio" gorm:"varchar(40);comment:个人简介"`
	Address    string    `json:"address" gorm:"varchar(120);comment:住址"`
	Location   string    `json:"location" gorm:"varchar(20);comment:经纬度（逗号分隔）"`
	Profession string    `json:"profession" gorm:"varchar(40);comment:职业"`
	Email      string    `json:"email" gorm:"varchar(20);comment:邮箱"`
	Gender     int64     `json:"gender" gorm:"comment:性别（0未选择|1男|2女|3武装直升机）"`
	Birthday   time.Time `json:"birthday" gorm:"default:null;comment:生日"`
}

// SimpleUserProfile 缩略用户信息结构体
type SimpleUserProfile struct {
	CreatorID uint   `json:"creatorId"`
	NickName  string `json:"nickName"`
	Avatar    string `json:"avatar"`
}

// TableName 设置简略用户信息结构体表名
func (SimpleUserProfile) TableName() string {
	return "user_profile"
}

// UserAuth 用户授权请求结构体
type UserAuth struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// UserIM im信息
type UserIM struct {
	IMUserId  string `json:"imUserId" gorm:"varchar(20);not null;unique;comment:im登录账号"`
	IMToken   string `json:"imToken" gorm:"-"`
	IMExpired int64  `json:"imExpired" gorm:"-"`
}
