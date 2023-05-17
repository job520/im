package logic

import (
	"encoding/json"
	"fmt"
	"im/websocket/service"
	"im/websocket/variables"
)

// 后端调度逻辑处理
func dispatch(data string) {
	msg := variables.Message{}
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("msg from client:", msg.Msg)
	switch msg.CMD {
	case variables.SingleMsg:
		service.ReceiveSingleMsg(msg.DestID, "hello from server!")
	case variables.GroupMsg:
		service.ReceiveGroupMsg(msg.DestID, "hello from server!")
	case variables.HeartMsg:
		// 检测客户端的心跳
	}
}
