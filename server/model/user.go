package model

import (
	"time"
)

// UserModel 用户结构体
type UserModel struct {
	OrmModel

	PhoneNumber string            `json:"phoneNumber" gorm:"varchar(40);not null;unique;comment:手机号（登录凭证）"`
	Password    string            `json:"-" gorm:"size:255;not null;comment:密码（登录密码）"`
	Profile     *UserProfileModel `json:"profile,omitempty" gorm:"-"`
}

// RespUserModel 报文用户结构体
type RespUserModel struct {
	ID int64 `json:"id" gorm:"primarykey"`

	Profile *RespUserProfileModel `json:"profile,omitempty" gorm:"-"`
}

// UserProfileModel 用户信息结构体
type UserProfileModel struct {
	OrmModel
	RespUserProfileModel

	Avatar       string          `json:"avatar" gorm:"varchar(40);comment:头像（只存储oss的key/id）"`
	Bio          string          `json:"bio" gorm:"varchar(40);comment:个人简介"`
	AddressCodes []string        `json:"addressCodes" gorm:"type:json;serializer:json;comment:住址字典码集合"`
	Address      []RespDictModel `json:"address" gorm:"-"`
	Location     string          `json:"location" gorm:"varchar(20);comment:经纬度（逗号分隔）"`
	Profession   string          `json:"profession" gorm:"varchar(40);comment:职业"`
	GenderCode   string          `json:"genderCode" gorm:"comment:性别字典码"`
	Birth        time.Time       `json:"birth" gorm:"default:null;comment:生日"`

	ShipAddress []UserShipAddressModel `json:"shipAddress" gorm:"-"`
	Config      UserConfigModel        `json:"config" gorm:"-"`
	Level       UserLevelModel         `json:"level" gorm:"-"`
}

// RespUserProfileModel 报文用户信息结构体
type RespUserProfileModel struct {
	CreatorModel

	NickName string        `json:"nickName" gorm:"varchar(20);comment:昵称"`
	Avatar   string        `json:"avatar" gorm:"varchar(40);comment:头像（只存储oss的key/id）"`
	MedalIds *[]int64      `json:"medalIds,omitempty" gorm:"type:json;serializer:json;comment:勋章id集合"`
	Medals   *[]MedalModel `json:"medals,omitempty" gorm:"-"`
}

// UserShipAddressModel 用户收货地址
type UserShipAddressModel struct {
	OrmModel
	CreatorModel

	Receiver      string          `json:"receiver" gorm:"varchar(20);not null;comment:收货人"`
	Contact       string          `json:"contact" gorm:"varchar(40)not null;comment:联系方式（手机号等）"`
	AddressCodes  []string        `json:"addressCodes" gorm:"type:json;serializer:json;comment:住址字典码集合"`
	Address       []RespDictModel `json:"address" gorm:"-"`
	AddressDetail string          `json:"addressDetail" gorm:"varchar(300);not null;comment:详细地址"`
	TagCode       string          `json:"tagCode" gorm:"comment:标签字典码"`
	Tag           *RespDictModel  `json:"tag" gorm:"-"`
	Default       bool            `json:"default" gorm:"not null;comment:是否为默认收货地址"`
	Order         int64           `json:"order" gorm:"not null;comment:排序"`
}

// UserConfigModel 用户偏好配置
type UserConfigModel struct {
	OrmModel
	CreatorModel

	EvaluateCode       string   `json:"evaluateCode" gorm:"not null;comment:自我评价字典码"`
	RecipeCuisineCodes []string `json:"recipeCuisineCodes" gorm:"type:json;serializer:json;comment:偏好食谱菜系字典码集合"`
	RecipeTasteCodes   []string `json:"recipeTasteCodes" gorm:"type:json;serializer:json;comment:偏好食谱口味字典码集合"`
}

// UserLevelModel 用户等级
type UserLevelModel struct {
	OrmModel
	CreatorModel

	Exp       int64 `json:"exp" gorm:"not null;comment:用户获得的经验"`
	Level     int64 `json:"level" gorm:"not null;comment:用户当前等级"`
	LevelExp  int64 `json:"levelExp" gorm:"not null;当前等级已获得经验"`
	UpdateExp int64 `json:"updateExp" gorm:"not null;升级所需经验"`
}
