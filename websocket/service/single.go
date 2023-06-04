package service

import "im/websocket/global"

// 发送消息,发送到消息的管道
func ReceiveSingleMsg(userId string, msg string) {
	global.RwLocker.RLock()
	node, ok := global.ClientMap[userId]
	global.RwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
