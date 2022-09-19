package common

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger 初始化日志系统
func InitLogger() *zap.Logger {
	logger := zap.NewExample()
	Logger = logger
	return logger
}

// LogInfo 记录日志
func LogInfo(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// GetLogger 获取日志对象
func GetLogger() *zap.Logger {
	return Logger
}
