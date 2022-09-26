package model

import (
	"time"
)

// User 用户结构体
type User struct {
	OrmBase
	SimpleUser
	UserConfig
	UserLevel

	PhoneNumber string    `json:"phoneNumber" gorm:"varchar(40);not null;unique;comment:手机号（登录凭证）"`
	Password    string    `json:"-" gorm:"varchar(120);not null;comment:密码（登录密码）"`
	Bio         string    `json:"bio" gorm:"varchar(40);comment:个人简介"`
	Location    string    `json:"location" gorm:"varchar(20);comment:经纬度（逗号分隔）"`
	Profession  string    `json:"profession" gorm:"varchar(40);comment:职业"`
	GenderCode  string    `json:"genderCode" gorm:"comment:性别字典码"`
	Birth       time.Time `json:"birth" gorm:"default:null;comment:生日"`
	MedalIds    []int64   `json:"medalIds" gorm:"type:json;serializer:json;comment:勋章id集合"`

	// 关注列表
	Subscribes []User `json:"-" gorm:"many2many:user_subscribes;comment:用户订阅列表"`
}

// SimpleUser 简易用户结构体
type SimpleUser struct {
	ID       int64  `json:"id" gorm:"primarykey"`
	NickName string `json:"nickName" gorm:"varchar(40);comment:昵称"`
	Avatar   string `json:"avatar" gorm:"varchar(80);comment:头像（只存储oss的key/id）"`
}

func (SimpleUser) TableName() string {
	return "user"
}

// UserConfig 用户偏好配置
type UserConfig struct {
	EvaluateCode       string   `json:"evaluateCode" gorm:"not null;comment:自我评价字典码"`
	RecipeCuisineCodes []string `json:"recipeCuisineCodes" gorm:"type:json;serializer:json;comment:偏好食谱菜系字典码集合"`
	RecipeTasteCodes   []string `json:"recipeTasteCodes" gorm:"type:json;serializer:json;comment:偏好食谱口味字典码集合"`
}

// UserLevel 用户等级
type UserLevel struct {
	Exp       int64 `json:"exp" gorm:"not null;comment:用户获得的经验"`
	Level     int64 `json:"level" gorm:"not null;comment:用户当前等级"`
	LevelExp  int64 `json:"levelExp" gorm:"not null;当前等级已获得经验"`
	UpdateExp int64 `json:"updateExp" gorm:"not null;升级所需经验"`
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
