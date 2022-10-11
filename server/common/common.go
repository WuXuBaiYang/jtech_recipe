package common

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/tool"
	"time"
)

// ServerHost 服务基本地址
var ServerHost = tool.If[string](gin.IsDebugging(), devServerHost, productionServerHost)

// 开发地址
const devServerHost = "127.0.0.1"

// 生产地址 api.jtech.live
const productionServerHost = ""

// SMSExpirationTime 短信验证码的有效时间
const SMSExpirationTime = 5 * time.Minute

// 日志配置信息
var logConfig = struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Skip       int
	FuncKey    string
	FileKey    string
	LineKey    string
}{
	Filename:   tool.If(gin.IsDebugging(), "C:/Users/wuxubaiyang/dev.log", "/tmp/jtech_server.log"),
	MaxSize:    100,
	MaxBackups: 60,
	MaxAge:     1,
	Compress:   false,
	Skip:       2,
	FuncKey:    "func",
	FileKey:    "file",
	LineKey:    "line",
}

// 数据库配置信息
var dbConfig = struct {
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
	UserName:  "root",
	Password:  "fmekST31BnzvPa",
	Host:      ServerHost,
	Port:      "3306",
	Database:  "jtech_recipe",
	Charset:   "utf8mb4",
	ParseTime: "True",
	Loc:       "Local",
}

// redis数据库配置
var rdbConfig = struct {
	Addr     string
	Port     int
	Password string
}{
	Addr:     ServerHost,
	Port:     6379,
	Password: "X0RrnycFiMt8ab",
}

// jwt授权配置信息
var jwtConfig = struct {
	Key                   []byte
	ExpirationTime        time.Duration
	RefreshExpirationTime time.Duration
	Issuer                string
}{
	Key:                   []byte("jtech_jh_server"),
	ExpirationTime:        tool.If(gin.IsDebugging(), 30*24*time.Hour, 15*time.Minute),
	RefreshExpirationTime: 30 * 24 * time.Hour,
	Issuer:                "jtech@127.0.0.1",
}
