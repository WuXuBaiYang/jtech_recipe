package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"server/model"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	username := "root"
	password := "JXuIAi4wqP0kho"
	host := ServerHost
	port := "13306"
	database := "jtech_im"
	charset := "utf8mb4"
	parseTime := "True"
	loc := "Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		username, password, host, port, database, charset, parseTime, loc)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	//迁移
	db.AutoMigrate(&model.User{}, &model.UserProfile{})
	db.AutoMigrate(&model.Post{}, &model.PostComment{}, &model.PostCommentReplay{}, &model.PostTag{})
	db.AutoMigrate(&model.Notification{})
	DB = db
	return db
}

// GetDB 获取数据库对象
func GetDB() *gorm.DB {
	return DB
}
