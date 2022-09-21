package model

import (
	"gorm.io/gorm"
	"time"
)

// Pagination 分页结构体
type Pagination struct {
	PageIndex   int         `json:"pageIndex"`
	PageSize    int         `json:"pageSize"`
	Total       int64       `json:"total"`
	CurrentSize int         `json:"currentSize"`
	Data        interface{} `json:"data"`
}

// AuthWithProfile 授权信息与用户信息结构体
type AuthWithProfile struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

// OrmModel gorm基类
type OrmModel struct {
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

// CreatorModel 创建者结构体
type CreatorModel struct {
	Creator   *SimpleUser `json:"creator,omitempty" gorm:"foreignKey:CreatorID;references:ID;"`
	CreatorID uint        `json:"creatorId" gorm:"comment:创建者id"`
}
