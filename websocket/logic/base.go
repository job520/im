package logic

import (
	"fmt"
	"github.com/gorilla/websocket"
	"im/websocket/config"
	"im/websocket/driver"
	"im/websocket/global"
	"im/websocket/service"
)

// 发送逻辑
func Send(node *global.Node) {
	for {
		select {
		case data := <-node.DataQueue: // 收到客户端消息
			err := node.Conn.WriteMessage(websocket.TextMessage, []byte(data)) // 服务端回应客户端消息
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}

// 接收逻辑
func Receive(node *global.Node) {
	for {
		_, data, err := node.Conn.ReadMessage() // 监听服务端接收到的消息
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		dispatch(string(data)) // 消息处理
	}
}

func CheckToken(token string) (string, bool) {
	userId, platform, isLegal := service.ParseJwtToken(config.Config.Jwt.EncryptKey, token)
	if !isLegal {
		return "", false
	}
	// 查看 redis 中有没有保存 token
	rdb := driver.NewRedisClient()
	redisKey := fmt.Sprintf("%s:%d", userId, platform)
	if res := rdb.Get(redisKey); res.Err() == nil {
		return userId, true
	} else {
		return "", false
	}
}
