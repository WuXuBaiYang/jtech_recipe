package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"server/common"
	"server/controller/response"
	"server/model"
	"strings"
)

// AuthCheck 授权校验
func AuthCheck(c *gin.Context) {
	claims, err := GetTokenClaim(c)
	if err != nil {
		response.FailAuthDef(c)
		c.Abort()
		return
	}
	// 判断当前平台与token中存储的平台是否保持一致
	if claims.Platform != GetPlatform(c) {
		response.FailParams(c, "禁止跨平台使用token")
		return
	}
	// 将用户信息写入上下文
	c.Set("userId", claims.UserId)
	c.Next()
}

// GetToken 获取token
func GetToken(c *gin.Context) string {
	return c.GetHeader("Authorization")
}

// GetTokenClaim 获取授权token信息
func GetTokenClaim(c *gin.Context) (*common.AccessClaims, error) {
	tokenString := GetToken(c)
	if len(tokenString) == 0 || !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.New("token不存在/格式错误")
	}
	token, claims, err := common.ParseAccessToken(tokenString[7:])
	if err != nil || !token.Valid {
		return nil, errors.New("token校验失败")
	}
	return claims, nil
}

// GetRefreshToken 获取refreshToken
func GetRefreshToken(c *gin.Context) string {
	return c.GetHeader("RefreshToken")
}

// GetRefreshTokenClaim 获取刷新token信息
func GetRefreshTokenClaim(c *gin.Context) (*common.RefreshClaims, error) {
	tokenString := GetRefreshToken(c)
	if len(tokenString) == 0 {
		return nil, errors.New("refreshToken不存在")
	}
	token, claims, err := common.ParseRefreshToken(tokenString)
	if err != nil || !token.Valid {
		return nil, errors.New("refreshToken校验失败")
	}
	return claims, nil
}

// GetCurrUId 获取当前用户id
func GetCurrUId(c *gin.Context) string {
	return c.GetString("userId")
}

// GetCurrUser 获取当前用户信息
func GetCurrUser(c *gin.Context) *model.User {
	uId := GetCurrUId(c)
	if len(uId) != 0 {
		db := common.GetDB()
		user := &model.User{}
		db.Find(&user, uId)
		return user
	}
	return nil
}

// PermissionCheck 权限校验
func PermissionCheck(c *gin.Context) {
	c.Next()
}
