package main

import (
	"database/sql"
	"fmt"
	"server/common"
)

func main() {
	//// 初始化日志系统
	//logger := common.InitLogger()
	//defer func(logger *zap.Logger) {
	//	err := logger.Sync()
	//	if err != nil {
	//		fmt.Printf("日志系统启动失败：%v", err)
	//		return
	//	}
	//}(logger)
	// 初始化数据库
	sqlDB, _ := common.InitDB().DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			fmt.Printf("数据库启动失败：%v", err)
		}
	}(sqlDB)
	//// 创建默认路由引擎
	//r := gin.Default()
	//// 注册路由
	//CollectRoutes(r)
	//// 启动服务
	//panic(r.Run(":9527"))
}
