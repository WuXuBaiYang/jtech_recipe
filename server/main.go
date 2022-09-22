package main

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"server/common"
)

var ctx = context.Background()

func main() {
	// 初始化日志系统
	logger := common.InitLogger()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic("日志系统启动失败" + err.Error())
		}
	}(logger)
	// 初始化数据库
	sqlDB, _ := common.InitDB().DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic("数据库关闭失败：" + err.Error())
		}
	}(sqlDB)
	// 初始化redis数据库
	rdb := common.InitRDB(ctx)
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			panic("redis数据库关闭失败" + err.Error())
		}
	}(rdb)
	// 创建默认路由引擎
	r := gin.Default()
	// 注册路由
	CollectRoutes(r)
	// 启动服务
	panic(r.Run(":9527"))
}
