package common

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"server/model"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	username := "root"
	password := "fmekST31BnzvPa"
	host := ServerHost
	port := "3306"
	database := "jtech_recipe"
	charset := "utf8mb4"
	parseTime := "True"
	loc := "Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		username, password, host, port, database, charset, parseTime, loc)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	//迁移
	//db.AutoMigrate(&model.Dict{})
	//db.Create(&model.Dict{
	//	PCode: 0,
	//	Tag:   "test",
	//	Order: 0,
	//	Code:  100,
	//	State: true,
	//	Desc:  "aaa",
	//})
	type DictActivityType struct {
		model.Dict
	}
	var a = &DictActivityType{}
	db.First(&a)
	//db.AutoMigrate(&model.User{}, &model.UserProfile{})
	//db.AutoMigrate(&model.Post{}, &model.PostComment{}, &model.PostCommentReplay{}, &model.PostTag{})
	//db.AutoMigrate(&model.Notification{})
	var f = a.UpdatedAt.String()
	println(f)
	DB = db
	return db
}

// GetDB 获取数据库对象
func GetDB() *gorm.DB {
	return DB
}
