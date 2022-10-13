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
	claims, err := GetAccessTokenClaim(c)
	if err != nil {
		response.FailAuthDef(c)
		c.Abort()
		return
	}
	// 将用户信息写入上下文
	c.Set("userId", claims.UserId)
	c.Next()
}

// PermissionCheck 权限校验
func PermissionCheck(permissions []model.PermissionLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := GetCurrUId(c)
		permission := common.GetUserPermission(c, userId)
		for _, it := range permissions {
			if it == permission {
				c.Next()
				return
			}
		}
		response.FailDef(c, -1, "无访问权限")
		c.Abort()
		return
	}
}

// GetAccessToken 获取token
func GetAccessToken(c *gin.Context) string {
	return c.GetHeader("Authorization")
}

// GetAccessTokenClaim 获取授权token信息
func GetAccessTokenClaim(c *gin.Context) (*common.AccessClaims, error) {
	tokenString := GetAccessToken(c)
	if len(tokenString) == 0 || !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.New("token不存在/格式错误")
	}
	token, claims, err := common.ParseAccessToken(c, tokenString[7:])
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
	token, claims, err := common.ParseRefreshToken(c, tokenString)
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
