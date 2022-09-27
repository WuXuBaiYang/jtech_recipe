package common

import (
	"github.com/golang-jwt/jwt/v4"
	"server/model"
	"time"
)

// AccessClaims 授权token中的claim
type AccessClaims struct {
	UserId   int64
	UserName string
	Platform string
	jwt.RegisteredClaims
}

// RefreshClaims 刷新token中的claim
type RefreshClaims struct {
	Target string
	jwt.RegisteredClaims
}

// ReleaseAccessToken 发放授权token
func ReleaseAccessToken(user model.User, platform string) (string, error) {
	expiresAt := time.Now().Add(jwtConfig.ExpirationTime)
	claims := &AccessClaims{
		UserId:   user.ID,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
			Subject:   user.PhoneNumber,
		},
	}
	return ReleaseToken(claims)
}

// ReleaseRefreshToken 发放刷新token
func ReleaseRefreshToken(user model.User, target string) (string, error) {
	expiresAt := time.Now().Add(jwtConfig.RefreshExpirationTime)
	claims := &RefreshClaims{
		Target: target,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
			Subject:   user.PhoneNumber,
		},
	}
	return ReleaseToken(claims)
}

// ParseAccessToken 授权token解析
func ParseAccessToken(tokenString string) (*jwt.Token, *AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtConfig.Key, nil
	})
	return token, claims, err
}

// ParseRefreshToken 刷新token解析
func ParseRefreshToken(tokenString string) (*jwt.Token, *RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtConfig.Key, nil
	})
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
