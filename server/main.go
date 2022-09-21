package main

import (
	"database/sql"
	"server/common"
)

func main() {
	// 初始化数据库
	sqlDB, _ := common.InitDB().DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("数据库关闭失败：" + err.Error())
		}
	}(sqlDB)
}
