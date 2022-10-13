package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"server/model"
	"time"
)

// 用户权限等级
const userPermissionKey = "user_permission"

// 在线用户存储标记
const onlineUserKey = "online_user"

// 授权存储key
const accessTokenKey = "access_token"

// 刷新存储key
const refreshTokenKey = "refresh_token"

// AccessClaims 授权token中的claim
type AccessClaims struct {
	UserId string
	jwt.RegisteredClaims
}

// RefreshClaims 刷新token中的claim
type RefreshClaims struct {
	UserId string
	Target string
	jwt.RegisteredClaims
}

// ReleaseAccessToken 发放授权token
func ReleaseAccessToken(c context.Context, user model.User) (string, error) {
	expiresAt := time.Now().Add(jwtConfig.ExpirationTime)
	claims := &AccessClaims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
			Subject:   fmt.Sprintf("%s", user.ID),
		},
	}
	token, err := ReleaseToken(claims)
	if len(token) != 0 {
		rdb := GetBaseRDB()
		rdb.Set(c, accessTokenKey+user.ID,
			token, jwtConfig.ExpirationTime)
		rdb.ZAdd(c, onlineUserKey, redis.Z{
			Score:  float64(expiresAt.UnixMilli()),
			Member: user.ID,
		})
		// 如果非普通用户登录，则记录用户权限等级
		if user.Permission != model.GeneralUser {
			rdb.Set(c, userPermissionKey+user.ID,
				int64(user.Permission), jwtConfig.ExpirationTime)
		}
	}
	return token, err
}

// ReleaseRefreshToken 发放刷新token
func ReleaseRefreshToken(c context.Context, user model.User, target string) (string, error) {
	expiresAt := time.Now().Add(jwtConfig.RefreshExpirationTime)
	claims := &RefreshClaims{
		UserId: user.ID,
		Target: target,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
			Subject:   fmt.Sprintf("%s", user.ID),
		},
	}
	token, err := ReleaseToken(claims)
	if len(token) != 0 {
		rdb := GetBaseRDB()
		rdb.Set(c, refreshTokenKey+user.ID,
			token, jwtConfig.RefreshExpirationTime)
	}
	return token, err
}

// ParseAccessToken 授权token解析
func ParseAccessToken(c context.Context, tokenString string) (*jwt.Token, *AccessClaims, error) {
	var claims AccessClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtConfig.Key, nil
		})
	if err == nil {
		rdb := GetBaseRDB()
		key := accessTokenKey + claims.UserId
		if n, err := rdb.Exists(c, key).
			Result(); err != nil || n <= 0 {
			return nil, nil, errors.New("token无效")
		}
	}
	return token, &claims, err
}

// ParseRefreshToken 刷新token解析
func ParseRefreshToken(c context.Context, tokenString string) (*jwt.Token, *RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtConfig.Key, nil
		})
	if err == nil {
		rdb := GetBaseRDB()
		key := refreshTokenKey + claims.UserId
		if n, err := rdb.Exists(c, key).
			Result(); err != nil || n <= 0 {
			return nil, nil, errors.New("刷新token无效")
		}
	}
	return token, claims, err
}

// ReleaseToken 发放token
func ReleaseToken(claims jwt.Claims) (string, error) {
	// jwt签名生成密钥
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtConfig.Key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ClearRDBToken 清除授权token和刷新token
func ClearRDBToken(c context.Context, ids ...string) *redis.IntCmd {
	rdb := GetBaseRDB()
	var keys []string
	for _, id := range ids {
		keys = append(keys, accessTokenKey+id,
			refreshTokenKey+id)
	}
	rdb.ZRem(c, onlineUserKey, ids)
	return rdb.Del(c, keys...)
}

// GetUserPermission 获取用户权限等级
func GetUserPermission(c context.Context, userId string) model.PermissionLevel {
	rdb := GetBaseRDB()
	cmd := rdb.Get(c, userPermissionKey+userId)
	if cmd.Err() != nil {
		return model.GeneralUser
	}
	v, _ := cmd.Int64()
	return model.PermissionLevel(v)
}
