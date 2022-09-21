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

	Avatar       string           `json:"avatar" gorm:"varchar(40);comment:头像（只存储oss的key/id）"`
	Bio          string           `json:"bio" gorm:"varchar(40);comment:个人简介"`
	AddressCodes []string         `json:"addressCodes" gorm:"type:json;serializer:json;comment:住址字典码集合"`
	Address      *[]RespDictModel `json:"address,omitempty" gorm:"-"`
	Location     string           `json:"location" gorm:"varchar(20);comment:经纬度（逗号分隔）"`
	Profession   string           `json:"profession" gorm:"varchar(40);comment:职业"`
	GenderCode   string           `json:"genderCode" gorm:"comment:性别字典码"`
	Birth        time.Time        `json:"birth" gorm:"default:null;comment:生日"`
}

// RespUserProfileModel 报文用户信息结构体
type RespUserProfileModel struct {
	CreatorModel

	NickName string `json:"nickName" gorm:"varchar(20);comment:昵称"`
	Avatar   string `json:"avatar" gorm:"varchar(40);comment:头像（只存储oss的key/id）"`
}
