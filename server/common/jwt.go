package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"server/model"
	"server/tool"
	"time"
)

// blockOutKey 账号封锁存储标记
const blockOutKey = "block_out_key"

// 授权存储key
const accessKey = "access_key"

// 刷新存储key
const refreshKey = "refresh_key"

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
		rdb.ZAdd(c, accessKey, redis.Z{
			Score:  float64(expiresAt.UnixMilli()),
			Member: tool.JoinV(user.ID, token),
		})
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
		rdb.ZAdd(c, refreshKey, redis.Z{
			Score:  float64(expiresAt.UnixMilli()),
			Member: tool.JoinV(user.ID, token),
		})
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
		rdb := GetBaseRDB()
		cmd := rdb.ZRank(c, accessKey,
			tool.JoinV(claims.UserId, tokenString))
		if cmd.Err() != nil && cmd.String() != tokenString {
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
		rdb := GetBaseRDB()
		cmd := rdb.ZRank(c, refreshKey,
			tool.JoinV(claims.UserId, tokenString))
		if cmd.Err() != nil && cmd.String() != tokenString {
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
	//return rdb.Del(c, keys...)
	return nil
}

// BlockOutUser 写入账号封锁记录
func BlockOutUser(c context.Context, ids ...string) *redis.IntCmd {
	rdb := GetBaseRDB()
	var members []redis.Z
	for _, id := range ids {
		members = append(members, redis.Z{Member: id})
	}
	return rdb.ZAdd(c, blockOutKey, members...)
}

// CheckBlockOut 检查账号是否已封锁
func CheckBlockOut(c context.Context, id string) bool {
	rdb := GetBaseRDB()
	cmd := rdb.ZRank(c, blockOutKey, id)
	return cmd.Err() == nil
}
