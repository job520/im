package logic

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"im/websocket/global"
	"log"
)

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := global.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	switch msg.Cmd {
	case global.CmdSingleMsg:
		SendMsg(msg.Dstid, data)
	case global.CmdRoomMsg:
		for _, v := range global.ClientMap {
			if v.GroupSets.Has(msg.Dstid) {
				v.DataQueue <- data
			}
		}
	case global.CmdHeart:
		// 检测客户端的心跳
	}
}

// 发送逻辑
func Send(node *global.Node) {
	for {
		select {
		case data := <-node.DataQueue: // 收到客户端消息
			err := node.Conn.WriteMessage(websocket.TextMessage, data) // 服务端回应客户端消息
			if err != nil {
				log.Println(err.Error())
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
			log.Println(err.Error())
			return
		}

		dispatch(data) // 消息处理

		fmt.Printf("recv<=%s", data)
	}
}

// 发送消息,发送到消息的管道
func SendMsg(userId int64, msg []byte) {
	global.RwLocker.RLock()
	node, ok := global.ClientMap[userId]
	global.RwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func CheckToken(userId int64, token string) bool {
	return true
}
