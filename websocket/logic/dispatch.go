package logic

import (
	"encoding/json"
	"fmt"
	"im/websocket/global"
	"im/websocket/service"
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
		service.ReceiveSingleMsg(msg.DestID, "hello from server!")
	case global.GroupMsg:
		service.ReceiveGroupMsg(msg.DestID, "hello from server!")
	case global.HeartMsg:
		// 检测客户端的心跳
	}
}
