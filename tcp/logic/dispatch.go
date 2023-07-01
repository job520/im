package logic

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"im/tcp/config"
	"im/tcp/global"
	"im/tcp/service"
)

// 后端调度逻辑处理
func Dispatch(userId string, platform int, msg global.Message) {
	msg.FromId = userId
	fmt.Println("msg from client:", msg.Msg)
	// 更新 userId:platform -> wsServer 状态
	userStatus := service.NewUserStatus(userId, platform, config.Config.Server.Address)
	if err := userStatus.Online(20); err != nil {
		logrus.Errorf("userStatus.Online error:%s\n", err.Error())
	}
	err := service.ReceiveSingleMsg(msg.ToId, platform, msg.Msg)
	if err != nil {
		logrus.Errorf("消息发送失败，error:%s\n", err.Error())
	}
}
