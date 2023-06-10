package service

import (
	"fmt"
	"im/websocket/global"
)

// 发送消息,发送到消息的管道
func ReceiveSingleMsg(userId string, platform int, msg string) {
	global.RwLocker.RLock()
	mapKey := fmt.Sprintf("%s:%d", userId, platform)
	node, ok := global.ClientMap[mapKey]
	global.RwLocker.RUnlock()
	if ok {
		node.MsgChan <- msg
	}
}
