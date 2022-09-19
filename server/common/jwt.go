package common

import (
	"github.com/golang-jwt/jwt/v4"
	"server/model"
	"time"
)

// jwt密钥
var jwtKey = []byte("jtech_jh_server")

// AccessClaims 授权token中的claim
type AccessClaims struct {
	UserId   uint
	UserName string
	Platform string
	jwt.RegisteredClaims
}

// RefreshClaims 刷新token中的claim
type RefreshClaims struct {
	Target string
	jwt.RegisteredClaims
}

// ReleaseToken 发放token
func ReleaseToken(claims jwt.Claims) (string, error) {
	// jwt签名生成密钥
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ReleaseAccessToken 发放授权token
func ReleaseAccessToken(user model.User, platform string) (string, error) {
	// token有效期
	//** 测试代码提供无限失效时间，正式上线后需要删除 **//
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	//expirationTime := time.Now().Add(15 * time.Minute)
	claims := &AccessClaims{
		UserId:   user.ID,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "127.0.0.1",
			Subject:   "user token",
		},
	}
	return ReleaseToken(claims)
}

// ReleaseRefreshToken 发放刷新token
func ReleaseRefreshToken(target string) (string, error) {
	// token有效期
	expirationTime := time.Now().Add(120 * 24 * time.Hour)
	claims := &RefreshClaims{
		Target: target,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "127.0.0.1",
			Subject:   "user refresh token",
		},
	}
	return ReleaseToken(claims)
}

// ParseAccessToken 授权token解析
func ParseAccessToken(tokenString string) (*jwt.Token, *AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

// ParseRefreshToken 刷新token解析
func ParseRefreshToken(tokenString string) (*jwt.Token, *RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
