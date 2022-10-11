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

// 授权存储key
const accessKey = "_access_key"

// AccessClaims 授权token中的claim
type AccessClaims struct {
	UserId string
	jwt.RegisteredClaims
}

// 刷新存储key
const refreshKey = "_refresh_key"

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
		rdb := GetAuthRDB()
		rdb.Set(c, user.ID+accessKey,
			token, jwtConfig.ExpirationTime)
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
		rdb := GetAuthRDB()
		rdb.Set(c, user.ID+refreshKey,
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
	if token != nil {
		rdb := GetAuthRDB()
		v := rdb.Get(c, claims.UserId+accessKey)
		if v != nil && v.String() != tokenString {
			ClearRDBToken(c, claims.UserId)
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
	if token != nil {
		rdb := GetAuthRDB()
		v := rdb.Get(c, claims.UserId+refreshKey)
		if v != nil && v.String() != tokenString {
			ClearRDBToken(c, claims.UserId)
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
	rdb := GetAuthRDB()
	var keys []string
	for _, id := range ids {
		keys = append(keys,
			id+accessKey, id+refreshKey)
	}
	return rdb.Del(c, keys...)
}
