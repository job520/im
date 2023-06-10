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
		case msg := <-node.MsgChan: // 收到客户端消息
			err := node.Conn.WriteMessage(websocket.TextMessage, []byte(msg)) // 服务端回应客户端消息
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}

// 接收逻辑
func Receive(userId string, platform int, node *global.Node) {
	for {
		_, data, err := node.Conn.ReadMessage() // 监听服务端接收到的消息
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		dispatch(userId, platform, string(data)) // 消息处理
	}
}

func CheckToken(token string) (string, int, bool) {
	userId, platform, isLegal := service.ParseJwtToken(config.Config.Jwt.EncryptKey, token)
	if !isLegal {
		return "", 0, false
	}
	// 查看 redis 中有没有保存 token
	rdb := driver.NewRedisClient()
	redisKey := fmt.Sprintf(global.JwtKey, userId, platform)
	if res := rdb.Get(redisKey); res.Err() == nil {
		return userId, platform, true
	} else {
		return "", platform, false
	}
}
