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

// lucene solr/es wukong 分词检索框架

func main() {
	// 设置调试状态
	gin.SetMode(gin.DebugMode)
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
	rdbList := common.InitRDB(ctx)
	defer func(clients []*redis.Client) {
		for _, it := range clients {
			err := it.Close()
			if err != nil {
				panic("redis数据库关闭失败" + err.Error())
			}
		}
	}(rdbList)
	// 创建默认路由引擎
	r := gin.Default()
	// 注册验证方法
	common.InitValidator()
	// 注册路由
	CollectRoutes(r)
	// 启动服务
	panic(r.Run(":9527"))
}
