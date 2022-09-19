package model

type Dict struct {
	OrmModel
	PCode     int64  `json:"pCode" gorm:"not null;comment:父级键值"`
	CreatorId uint   `json:"creatorId" gorm:"comment:创建者id"`
	Tag       string `json:"tag" gorm:"varchar(40);not null; comment:标签"`
	Order     int64  `json:"order" gorm:"AUTO_INCREMENT;not null;comment:排序"`
	Code      int64  `json:"code" gorm:"not null;unique;comment:键值"`
	State     bool   `json:"state" gorm:"not null;comment:是否可用"`
	Desc      string `json:"desc" gorm:"varchar(80);comment:标记/描述"`
}
