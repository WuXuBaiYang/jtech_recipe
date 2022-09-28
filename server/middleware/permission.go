package middleware

import (
	"github.com/gin-gonic/gin"
)

// PermissionCheck 权限校验
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
