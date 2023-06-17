package service

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"im/websocket/global"
)

func ReceiveSingleMsg(userId string, platform int, msg string) error {
	global.ConnectMap.Lock()
	mapKey := fmt.Sprintf("%s:%d", userId, platform)
	conn, ok := global.ConnectMap.ClientMap[mapKey]
	if ok {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			return err
		}
	} else {
		return errors.New("未找到连接句柄")
	}
	global.ConnectMap.Unlock()
	return nil
}
