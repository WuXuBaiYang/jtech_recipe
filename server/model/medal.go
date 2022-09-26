package model

// Medal 勋章结构体
type Medal struct {
	OrmBase

	Logo       string `json:"logo" gorm:"not null;comment:勋章图标（oss的key/id）"`
	Name       string `json:"name" gorm:"not null;comment:勋章名称"`
	RarityCode string `json:"rarityCode" gorm:"not null;comment:勋章稀有度字典码"`
}
