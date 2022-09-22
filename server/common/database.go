package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"server/model"
)

var db *gorm.DB

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
	c := DBConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.UserName, c.ParseTime, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc)
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

// GetDB 获取数据库对象
func GetDB() *gorm.DB {
	return db
}
