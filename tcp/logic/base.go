package logic

import (
	"fmt"
	"im/tcp/config"
	"im/tcp/driver"
	"im/tcp/global"
	"im/tcp/service"
)

func CheckToken(token string) (string, int, bool) {
	userId, platform, isLegal := service.ParseJwtToken(config.Config.Jwt.EncryptKey, token)
	if !isLegal {
		return "", 0, false
	}
	// 查看 redis 中有没有保存 token
	rdb := driver.NewRedisClient()
	redisKey := fmt.Sprintf(global.RedisJwtKey, userId, platform)
	if res := rdb.Get(redisKey); res.Err() == nil {
		return userId, platform, true
	} else {
		return "", platform, false
	}
}
