package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"runtime"
)

var logger *zap.Logger

// InitLogger 日志初始化
func InitLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	fileWriteSyncer := getFileLogWriter()
	tee := []zapcore.Core{
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.InfoLevel),
	}
	core := zapcore.NewTee(tee...)
	logger = zap.New(core)
	return logger
}

// 日志分隔
func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
	c := logConfig
	lumberJackLogger := &lumberjack.Logger{
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// LogInfo info日志
func LogInfo(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Info(message, fields...)
}

// LogDebug debug日志
func LogDebug(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Debug(message, fields...)
}

// LogError error日志
func LogError(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Error(message, fields...)
}

// LogWarn warn日志
func LogWarn(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Warn(message, fields...)
}

// 获取日志格式
func getCallerInfoForLog() (callerFields []zap.Field) {
	c := logConfig
	pc, file, line, ok := runtime.Caller(c.Skip) // 回溯两层，拿到写日志的业务函数的信息
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base函数返回路径的最后一个元素，只保留函数名
	callerFields = append(callerFields, zap.String(c.FuncKey, funcName), zap.String(c.FileKey, file), zap.Int(c.LineKey, line))
	return
}
