package model

import (
	"gorm.io/gorm"
	"server/tool"
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

// AuthWithUser 授权信息与用户信息结构体
type AuthWithUser struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	User         UserModel `json:"user"`
}

// OrmModel gorm基类
type OrmModel struct {
	ID        int64          `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

// CreatorModel 创建者结构体
type CreatorModel struct {
	CreatorId int64      `json:"creatorId" gorm:"not null;comment:创建者id"`
	Creator   *UserModel `json:"creator,omitempty" gorm:"-"`
}

// CreateOrmModel 创建基础结构体
func CreateOrmModel() *OrmModel {
	if id := tool.GenID(); id != 0 {
		return &OrmModel{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return nil
}
