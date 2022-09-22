package common

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"server/tool"
	"strings"
)

// DebugMode 调试状态
const DebugMode = true

// ServerHost 服务基本地址
var ServerHost = tool.If[string](DebugMode, devServerHost, productionServerHost)

// 开发地址
const devServerHost = "127.0.0.1"

// 生产地址 api.jtech.live
const productionServerHost = ""

// DBConfig 数据库配置信息
var DBConfig = struct {
	UserName   string
	Password   string
	Host       string
	Port       string
	Database   string
	Charset    string
	ParseTime  string
	Loc        string
	GormConfig *gorm.Config
}{
	UserName:  "jtech_server",
	Password:  "JXuIAi4wqP0kho",
	Host:      ServerHost,
	Port:      "3306",
	Database:  "jtech_recipe",
	Charset:   "utf8mb4",
	ParseTime: "True",
	Loc:       "Local",
	GormConfig: &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "jtech_",
			NameReplacer:  strings.NewReplacer("Resp", "", "Model", ""),
			SingularTable: true,
		},
	},
}
