package model

// Dict 字典项结构体
type Dict struct {
	OrmBase
	Creator
	SimpleDict
}

// SimpleDict 字典项报文结构体
type SimpleDict struct {
	PCode string `json:"pCode" gorm:"not null;comment:父值"`
	Code  string `json:"code" gorm:"not null;unique;comment:值"`
	Tag   string `json:"tag" gorm:"not null;comment:标签"`
	Info  string `json:"info" gorm:"comment:字典存储信息"`
}
