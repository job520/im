package logic

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"im/websocket/global"
)

// 后端调度逻辑处理
func dispatch(data string) {
	msg := global.Message{}
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("msg from client:", msg.Msg)
	switch msg.CMD {
	case global.SingleMsg:
		SendMsg(msg.DestID, "hello from server!")
	case global.GroupMsg:
		for _, v := range global.ClientMap {
			if v.GroupSets.Has(msg.DestID) {
				v.DataQueue <- data
			}
		}
	case global.HeartMsg:
		// 检测客户端的心跳
	}
}

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

// 发送消息,发送到消息的管道
func SendMsg(userId int, msg string) {
	global.RwLocker.RLock()
	node, ok := global.ClientMap[userId]
	global.RwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func CheckToken(userId int, token string) bool {
	return true
}
