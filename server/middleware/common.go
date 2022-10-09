package middleware

import (
	"github.com/gin-gonic/gin"
)

// Common 通用方法
func Common(c *gin.Context) {
	c.Next()
}
