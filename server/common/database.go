package common

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"server/model"
	"server/tool"
)

// 要合并的表单
var dst = []any{
	// 用户
	&model.UserModel{},
	&model.UserProfileModel{},
	&model.UserShipAddressModel{},
	&model.UserConfigModel{},
	&model.UserLevelModel{},
	// 帖子
	&model.PostModel{},
	&model.PostCommentModel{},
	&model.PostCommentReplayModel{},
	// 消息通知
	&model.NotifyModel{},
	// 活动
	&model.ActivityModel{},
	&model.ActivityRecordModel{},
	// 成就
	&model.MedalModel{},
	// 食谱
	&model.RecipeModel{},
	&model.RecipeStepModel{},
	// 菜单
	&model.MenuMode{},
}

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	// 初始化雪花算法
	if tool.GenID() == 0 {
		panic("雪花算法初始化失败")
	}
	c := dbConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.UserName, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc)
	DB, err := gorm.Open(mysql.Open(dsn), c.GormConfig)
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}
	if err := DB.AutoMigrate(dst...); err != nil {
		panic("数据库自动合并失败：" + err.Error())
	}
	db = DB
	return db
}

var db *gorm.DB

// GetDB 获取数据库对象
func GetDB() *gorm.DB {
	return db
}

// InitRDB 初始化redis数据库
func InitRDB(ctx context.Context) *redis.Client {
	c := rdbConfig
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Addr, c.Port),
		Password: c.Password,
		DB:       c.DB,
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic("redis数据库初始化失败：" + err.Error())
	}
	return rdb
}

var rdb *redis.Client

// GetRDB 获取redis默认数据库
func GetRDB() *redis.Client {
	return rdb
}
