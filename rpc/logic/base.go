package logic

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"im/rpc/generate/transfer"
	"im/rpc/global"
	"io"
)

type TransferService struct {
	requests []*transfer.ChatRequestAndResponse
}

func (t TransferService) Chat(conn transfer.Transfer_ChatServer) error {
	ctx := conn.Context()
	mdClient, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("get metadata error")
	}
	if fromConnectors := mdClient.Get("fromConnector"); len(fromConnectors) > 0 {
		fmt.Printf("metadata from client - fromConnector: %s\n", fromConnectors[0])
		global.ConnectMap.Lock()
		global.ConnectMap.ClientMap[fromConnectors[0]] = conn
		global.ConnectMap.Unlock()
	}
	for {
		msg, err := conn.Recv()
		if err == io.EOF {
			fmt.Println("receive done!")
			return nil
		}
		if err != nil {
			fmt.Println("receive error:", err)
			return err
		}
		fmt.Printf("msg received: %v .\n", msg)
		dispatch(msg)
	}
}

func dispatch(msg *transfer.ChatRequestAndResponse) {
	// todo: 消息分发逻辑（心跳/转发）
	switch msg.MsgType {
	// 消息转发
	case int32(global.RpcMsgTypeTransfer):
		toConnector := msg.ToConnector
		conn, ok := global.ConnectMap.ClientMap[toConnector]
		if ok {
			if err := conn.Send(msg); err != nil {
				logrus.Errorf("转发消息失败：%s\n", err)
			}
		}
	}
}
