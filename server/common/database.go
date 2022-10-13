package common

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"server/model"
	"server/tool"
)

// 要合并的表单
var dst = []any{
	// 用户
	&model.User{},
	&model.UserMedal{},
	&model.UserAddress{},
	// 评论回复
	&model.Comment{},
	&model.Replay{},
	// 帖子
	&model.Post{},
	// 消息通知
	&model.Notify{},
	// 活动
	&model.Activity{},
	&model.ActivityRecord{},
	// 食谱
	&model.Recipe{},
	// 菜单
	&model.Menu{},
}

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	// 初始化雪花算法
	if len(tool.GenID()) == 0 {
		panic("雪花算法初始化失败")
	}
	c := dbConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.UserName, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc)
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		//Logger:                                   logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "sys_",
			SingularTable: true,
		},
	})
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
func InitRDB(ctx context.Context) []*redis.Client {
	rdbList := []*redis.Client{
		GetBaseRDB(),
	}
	for _, it := range rdbList {
		if _, err := it.Ping(ctx).Result(); err != nil {
			panic("redis数据库初始化失败：" + err.Error())
		}
	}
	return rdbList
}

// 缓存redis数据库
var rdbClientMap = map[int]*redis.Client{}

// GetRDB 根据下标获取redis数据库
func getRDB(i int) *redis.Client {
	client, ok := rdbClientMap[i]
	if !ok {
		cfg := rdbConfig
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
			Password: cfg.Password,
			DB:       i,
		})
		rdbClientMap[i] = client
	}
	return client
}

// GetBaseRDB 获取redis数据库
func GetBaseRDB() *redis.Client {
	return getRDB(0)
}
