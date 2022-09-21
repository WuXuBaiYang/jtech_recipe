package model

type Dict struct {
	OrmModel
	CreatorId *string `json:"creatorId" gorm:"comment:创建者id"`
	State     bool    `json:"state" gorm:"not null;comment:是否可用"`
	PCode     string  `json:"pCode" gorm:"not null;comment:父级键值"`
	Code      string  `json:"code" gorm:"not null;unique;comment:键值"`
	Tag       string  `json:"tag" gorm:"varchar(40);not null; comment:标签"`
	Desc      string  `json:"desc" gorm:"varchar(80);comment:标记/描述"`
}
