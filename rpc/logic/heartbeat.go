package logic

import (
	"fmt"
	"im/rpc/generate/transfer"
	"im/rpc/global"
	"time"
)

func HeartBeat() {
	for {
		for connectorId, conn := range global.ConnectMap.ClientMap {
			go func(connectorId string, conn transfer.Transfer_ChatServer) {
				request := &transfer.ChatRequestAndResponse{
					FromConnector: "xx:80",
					ToConnector:   connectorId,
					MsgType:       1,
					Message:       "pong!",
					Data: &transfer.Data{
						Id:         0,
						Cmd:        0,
						FromId:     "",
						DestId:     "",
						Msg:        "",
						MsgType:    0,
						AckMsgType: 0,
					},
				}
				if err := conn.Send(request); err != nil {
					fmt.Println("send err:", err)
					global.ConnectMap.Lock()
					delete(global.ConnectMap.ClientMap, connectorId)
					global.ConnectMap.Unlock()
				}
			}(connectorId, conn)
		}
		time.Sleep(3 * time.Second)
	}
}
