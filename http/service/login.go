package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im/http/config"
	"im/http/driver"
	"im/http/model"
	"time"
)

func Register(ctx *gin.Context, username, password string) (bool, error) {
	return model.Register(ctx, username, password)
}

func Login(ctx *gin.Context, username, password string, platform int) (string, error) {
	userId, err := model.GetUserId(ctx, username, password)
	if err != nil {
		return "", err
	}
	encryptKey := config.Config.Jwt.EncryptKey
	expireHours := config.Config.Jwt.ExpireHours
	// 生成 jwt-token
	token, err := GenerateJwtToken(encryptKey, expireHours, userId, platform)
	// jwt-token 存入 redis
	rdb := driver.NewRedisClient()
	redisKey := fmt.Sprintf("%s:%d", userId, platform)
	if err := rdb.Set(redisKey, token, time.Duration(expireHours)*time.Hour).Err(); err != nil {
		return "", err
	}
	return token, err
}
