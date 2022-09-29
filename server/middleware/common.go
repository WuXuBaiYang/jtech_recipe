package middleware

import (
	"github.com/gin-gonic/gin"
	"server/controller/response"
	"server/tool"
)

// Common 通用方法
func Common(c *gin.Context) {
	// 获取平台信息
	platform := GetPlatform(c)
	if !tool.PlatformVerify(platform) {
		response.FailAuth(c, "缺少平台信息")
		c.Abort()
		return
	}
	c.Next()
}

// GetPlatform 获取平台信息
func GetPlatform(c *gin.Context) string {
	return c.GetHeader("Platform")
}
