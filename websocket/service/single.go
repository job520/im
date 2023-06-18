package service

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"im/websocket/config"
	"im/websocket/generate/transfer"
	"im/websocket/global"
)

func ReceiveSingleMsg(userId string, platform int, msg string) error {
	global.ConnectMap.Lock()
	mapKey := fmt.Sprintf("%s:%d", userId, platform)
	conn, ok := global.ConnectMap.ClientMap[mapKey]
	if ok {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			return err
		}
	} else {
		// 通过状态管理器(redis)获取连接的 websocket服务器地址
		us := NewUserStatus(userId, platform, "")
		wsServer, err := us.GetOnlineUserServer()
		if err != nil {
			// todo: 给离线用户发送离线消息
			return errors.New("用户已离线")
		}
		// 通过 rpc转发
		global.RpcMsgChan <- &transfer.ChatRequestAndResponse{
			FromConnector: config.Config.Server.Address,
			ToConnector:   wsServer,
			MsgType:       int32(global.RpcMsgTypeTransfer),
			Message:       "",
			Data: &transfer.Data{
				Id:         0,
				Cmd:        int32(global.SingleMsg),
				FromId:     "",
				DestId:     "",
				Msg:        "a test message...",
				MsgType:    0,
				AckMsgType: 0,
			},
		}
	}
	global.ConnectMap.Unlock()
	return nil
}
