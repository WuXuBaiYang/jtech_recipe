package main

import (
	"server/common"
	"server/script"
)

func main() {
	// 初始化数据库
	db := common.InitDB()
	// 初始化字典表
	script.InitDict(db)
}
