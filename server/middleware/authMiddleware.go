package middleware

import (
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
	"strings"
)

// AuthMiddleware 授权中间件
func AuthMiddleware(tokenCheck bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		tokenString := c.GetHeader("Authorization")
		// 校验token结构完整
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			response.FailAuthDef(c)
			c.Abort()
			return
		}
		// 提取token有效部分
		tokenString = tokenString[7:]
		token, claims, err := common.ParseAccessToken(tokenString)
		if tokenCheck && (err != nil || !token.Valid) {
			response.FailAuthDef(c)
			c.Abort()
			return
		}
		// 判断当前平台与token中存储的平台是否保持一致
		platform, _ := c.Get("platform")
		if claims.Platform != platform {
			response.FailAuth(c, "禁止跨平台使用token")
			c.Abort()
			return
		}
		// 提取token中的用户id
		userId := claims.UserId
		db := common.GetDB()
		var user model.User
		db.First(&user, userId)
		// 判断用户是否存在
		if user.ID == 0 {
			response.FailAuthDef(c)
			c.Abort()
			return
		}
		// 将用户信息写入上下文
		c.Set("user", user)
		c.Next()
	}
}
