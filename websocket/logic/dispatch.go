package logic

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"im/websocket/config"
	"im/websocket/global"
	"im/websocket/service"
)

// 后端调度逻辑处理
func dispatch(userId string, platform int, data string) {
	msg := global.Message{}
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	msg.FromID = userId
	fmt.Println("msg from client:", msg.Msg)
	switch msg.CMD {
	case global.HeartMsg:
		// 更新 userId:platform -> wsServer 状态
		userStatus := service.NewUserStatus(userId, platform, config.Config.Server.Address)
		if err := userStatus.Online(20); err != nil {
			logrus.Errorf("ping error:%s\n", err.Error())
		}
		service.ReceiveSingleMsg(userId, platform, "pong!")
	case global.SingleMsg:
		service.ReceiveSingleMsg(msg.DestID, platform, "hello from server!")
	case global.GroupMsg:
		service.ReceiveGroupMsg(msg.DestID, "hello from server!")
	}
}
