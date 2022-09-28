package middleware

import (
	"github.com/gin-gonic/gin"
	"server/controller/response"
	"server/tool"
)

// Common 通用方法
func Common() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取平台信息
		platform := c.GetHeader("Platform")
		if !tool.PlatformVerify(platform) {
			response.FailAuth(c, "缺少平台信息")
			c.Abort()
			return
		}
		// 将信息写入上下文
		c.Set("platform", platform)
		c.Next()
	}
}
