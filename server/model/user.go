package model

import "time"

// User 用户结构体
type User struct {
	OrmBase
	// 基础信息/敏感内容
	PhoneNumber string `json:"phoneNumber" gorm:"varchar(40);not null;unique;comment:手机号（登录凭证）"`
	Password    string `json:"-" gorm:"varchar(120);not null;comment:密码（登录密码）"`
	Blocked     bool   `json:"-" gorm:"comment:是否被封锁"`
	// 用户详细信息
	NickName   string      `json:"nickName" gorm:"varchar(40);comment:昵称"`
	Avatar     string      `json:"avatar" gorm:"varchar(80);comment:头像（只存储oss的key/id）"`
	Bio        string      `json:"bio" gorm:"varchar(40);comment:个人简介"`
	Profession string      `json:"profession" gorm:"varchar(40);comment:职业"`
	GenderCode string      `json:"genderCode" gorm:"comment:性别字典码"`
	Birth      *time.Time  `json:"birth" gorm:"comment:生日"`
	Medals     []UserMedal `json:"medals,omitempty" gorm:"many2many:user_has_medals;comment:已获得的勋章列表"`
	// 用户配置相关
	EvaluateCode       string   `json:"evaluateCode" gorm:"not null;comment:自我评价字典码"`
	RecipeCuisineCodes []string `json:"recipeCuisineCodes" gorm:"type:json;serializer:json;comment:偏好食谱菜系字典码集合"`
	RecipeTasteCodes   []string `json:"recipeTasteCodes" gorm:"type:json;serializer:json;comment:偏好食谱口味字典码集合"`
	// 用户经验
	Exp       int64 `json:"exp" gorm:"not null;comment:用户获得的经验"`
	Level     int64 `json:"level" gorm:"not null;comment:用户当前等级"`
	LevelExp  int64 `json:"levelExp" gorm:"not null;当前等级已获得经验"`
	UpdateExp int64 `json:"updateExp" gorm:"not null;升级所需经验"`
	// 关注列表
	Subscribes []User `json:"-" gorm:"many2many:user_subscribes;comment:用户订阅列表"`
	// 点赞/收藏过的帖子
	LikePosts    []Post `json:"-" gorm:"many2many:post_like_users"`
	CollectPosts []Post `json:"-" gorm:"many2many:post_collect_users"`
	// 点赞/收藏过的菜单
	LikeMenus    []Menu `json:"-" gorm:"many2many:menu_like_users"`
	CollectMenus []Menu `json:"-" gorm:"many2many:menu_collect_users"`
	// 点赞/收藏过的食谱
	LikeRecipes    []Recipe `json:"-" gorm:"many2many:recipe_like_users"`
	CollectRecipes []Recipe `json:"-" gorm:"many2many:recipe_collect_users"`
}

// SimpleUser 简易用户结构体
type SimpleUser struct {
	ID       string `json:"id"`
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
	Level    int64  `json:"level"`
}

func (SimpleUser) TableName() string {
	return "sys_user"
}

// UserAddress 用户收货地址
type UserAddress struct {
	OrmBase
	Creator

	Receiver      string      `json:"receiver" gorm:"varchar(20);not null;comment:收货人"`
	Contact       string      `json:"contact" gorm:"varchar(40)not null;comment:联系方式（手机号等）"`
	AddressCodes  []string    `json:"addressCodes" gorm:"type:json;serializer:json;comment:住址字典码集合"`
	AddressDetail string      `json:"addressDetail" gorm:"varchar(300);not null;comment:详细地址"`
	TagCode       string      `json:"tagCode" gorm:"comment:标签字典码"`
	Tag           *SimpleDict `json:"tag" gorm:"-"`
	Default       bool        `json:"default" gorm:"not null;comment:是否为默认收货地址"`
	Order         int64       `json:"order" gorm:"not null;comment:排序"`
}

// UserMedal 用户勋章结构体
type UserMedal struct {
	OrmBase

	Logo       string `json:"logo" gorm:"unique;not null;comment:勋章图标（oss的key/id）"`
	Name       string `json:"name" gorm:"unique;not null;comment:勋章名称"`
	RarityCode string `json:"rarityCode" gorm:"not null;comment:勋章稀有度字典码"`
}
