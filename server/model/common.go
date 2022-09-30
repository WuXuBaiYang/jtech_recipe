package model

import (
	"gorm.io/gorm"
	"time"
)

// Pagination 分页结构体
type Pagination[T interface{}] struct {
	PageIndex int   `json:"pageIndex" form:"pageIndex" validate:"gte=1"`
	PageSize  int   `json:"pageSize" form:"pageSize" validate:"gte=10"`
	Total     int64 `json:"total"`
	Data      []T   `json:"data"`
}

// OrmBase gorm基类
type OrmBase struct {
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

// Creator 创建者结构体
type Creator struct {
	CreatorId string      `json:"creatorId" gorm:"comment:创建者id"`
	Creator   *SimpleUser `json:"creator,omitempty" gorm:"foreignKey:CreatorId"`
}
