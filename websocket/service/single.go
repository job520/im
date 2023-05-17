package service

import "im/websocket/variables"

// 发送消息,发送到消息的管道
func ReceiveSingleMsg(userId int, msg string) {
	variables.RwLocker.RLock()
	node, ok := variables.ClientMap[userId]
	variables.RwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
