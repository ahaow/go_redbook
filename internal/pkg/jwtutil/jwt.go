package jwtutil

import (
	"fmt"
	"time"

	"go_redbook/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 是项目自己的 JWT 载荷。
// RegisteredClaims 里包含过期时间、签发人、签发时间这些标准字段。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 根据用户信息生成登录 token。
func GenerateToken(cfg config.JwtConfig, userID uint, username string) (string, error) {
	now := time.Now()
	expires := now.Add(time.Duration(cfg.Expires) * time.Hour)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Key))
}

// ParseToken 校验 token 并返回 Claims。
// 后面写登录鉴权中间件时，可以直接复用这个函数。
func ParseToken(cfg config.JwtConfig, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
