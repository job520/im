package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"im/websocket/driver"
	"im/websocket/global"
	"time"
)

type userStatus struct {
	Rdb      *redis.Client
	RedisKey string
	RedisVal string
}

func NewUserStatus(userId string, platform int, wsServer string) userStatus {
	rdb := driver.NewRedisClient()
	redisKey := fmt.Sprintf(global.RedisStatusKey, userId, platform)
	userStatusObj := userStatus{
		Rdb:      rdb,
		RedisKey: redisKey,
		RedisVal: wsServer,
	}
	return userStatusObj
}

// 用户上线
func (u userStatus) Online(timeout int) error {
	err := u.Rdb.Set(u.RedisKey, u.RedisVal, time.Duration(timeout)*time.Second).Err()
	return err
}

// 用户下线
func (u userStatus) Offline() error {
	err := u.Rdb.Del(u.RedisKey).Err()
	return err
}

// 获取在线用户所连接的 websocket 服务器
func (u userStatus) GetOnlineUserServer() (string, error) {
	if res := u.Rdb.Get(u.RedisKey); res.Err() == nil {
		return res.Val(), nil
	} else {
		return "", res.Err()
	}
}
