package model

// DictModel 字典项结构体
type DictModel struct {
	OrmModel
	CreatorModel
	RespDictModel

	State bool `json:"state"`
}

// RespDictModel 字典项报文结构体
type RespDictModel struct {
	PCode string `json:"pCode"`
	Code  string `json:"code"`
	Tag   string `json:"tag"`
	Info  string `json:"info"`
}
