package service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"im/tcp/config"
	"im/tcp/generate/transfer"
	"im/tcp/global"
)

func ReceiveSingleMsg(userId string, platform int, msg string) error {
	logrus.Info("---------------------")
	logrus.Infof("userId:%s,platform:%d,msg:%s\n", userId, platform, msg)
	mapKey := fmt.Sprintf("%s:%d", userId, platform)
	conn, ok := global.ConnectMap.ClientMap[mapKey]
	if ok {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			return err
		}
	} else {
		// 通过状态管理器(redis)获取连接的 websocket服务器地址
		us := NewUserStatus(userId, platform, "")
		wsServer, err := us.GetOnlineUserServer()
		logrus.Infof("wsServer: %s\n", wsServer)
		if err != nil {
			// todo: 给离线用户发送离线消息
			return errors.New("用户已离线")
		}
		// 通过 rpc转发
		transferMsg := &transfer.ChatRequestAndResponse{
			FromConnector: config.Config.Server.Address,
			ToConnector:   wsServer,
			Data: &transfer.Data{
				FromId: "",
				ToId:   userId,
				Msg:    msg,
			},
		}
		logrus.Infof("转发消息中..., fromConnector: %s, toConnector: %s \n", transferMsg.FromConnector, transferMsg.ToConnector)
		global.RpcMsgChan <- transferMsg
	}
	return nil
}
