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
func Dispatch(userId string, platform int, data string) {
	msg := global.Message{}
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	msg.FromId = userId
	fmt.Println("msg from client:", msg.Msg)
	// 更新 userId:platform -> wsServer 状态
	userStatus := service.NewUserStatus(userId, platform, config.Config.Server.Address)
	if err := userStatus.Online(20); err != nil {
		logrus.Errorf("userStatus.Online error:%s\n", err.Error())
	}
	err = service.ReceiveSingleMsg(msg.ToId, platform, "hello from server!")
	if err != nil {
		logrus.Errorf("消息发送失败，error:%s\n", err.Error())
	}
}
