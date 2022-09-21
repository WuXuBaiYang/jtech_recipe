package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"server/model"
	"strings"
)

var db *gorm.DB

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	username := "jtech_server"
	password := "JXuIAi4wqP0kho"
	host := ServerHost
	port := "3306"
	database := "jtech_recipe"
	charset := "utf8mb4"
	parseTime := "True"
	loc := "Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		username, password, host, port, database, charset, parseTime, loc)
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "jtech_",
			NameReplacer:  strings.NewReplacer("Resp", "", "Model", ""),
			SingularTable: true,
		},
	})
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}
	dst := []any{
		&model.UserModel{},
		&model.UserProfileModel{},
		&model.NotifyModel{},
		&model.PostModel{},
		&model.PostCommentModel{},
		&model.PostCommentReplayModel{},
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
