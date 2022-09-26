package model

// Dict 字典项结构体
type Dict struct {
	OrmBase
	Creator
	SimpleDict

	State bool `json:"state"`
}

// SimpleDict 字典项报文结构体
type SimpleDict struct {
	PCode string `json:"pCode"`
	Code  string `json:"code"`
	Tag   string `json:"tag"`
	Info  string `json:"info"`
}
