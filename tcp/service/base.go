package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// 生成 jwt-token
func GenerateJwtToken(encryptKey string, expireHours int, uid string, platform int) (string, error) {
	expire := time.Duration(expireHours)
	exp := time.Now().Add(time.Hour * expire).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      uid,
		"platform": platform,
		"exp":      exp,
	})
	tokenString, err := token.SignedString([]byte(encryptKey))
	return tokenString, err
}

// 解析 jwt-token
func ParseJwtToken(encryptKey, tokenString string) (string, int, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(encryptKey), nil
	})
	if err != nil {
		return "", 0, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if uid, ok := claims["uid"].(string); ok {
			return uid, int(claims["platform"].(float64)), true
		} else {
			return "", 0, false
		}
	} else {
		return "", 0, false
	}
}
