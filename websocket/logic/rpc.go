package logic

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"im/websocket/config"
	"im/websocket/generate/transfer"
	"im/websocket/global"
	"io"
)

func RpcClient() {
	address := config.Config.RpcServer.Address
	logrus.Info("rpc server address:", address)
	srv, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer srv.Close()
	mdClient := metadata.Pairs(
		"fromConnector", config.Config.Server.Address,
	)
	ctx := metadata.NewOutgoingContext(context.Background(), mdClient)
	c := transfer.NewTransferClient(srv)
	conn, err := c.Chat(ctx)
	if err != nil {
		logrus.Error("create stream error:", err)
	}

	// 发送消息
	go func() {
		for {
			select {
			case msg := <-global.RpcMsgChan:
				if err := conn.Send(msg); err != nil {
					return
				}
			}
		}
	}()

	// 接收消息
	for {
		msg, err := conn.Recv()
		if err == io.EOF {
			logrus.Info("receive done!")
			break
		}
		if err != nil {
			logrus.Error("receive error:", err)
		}
		logrus.Info("msg from server:", msg.Message)
		switch msg.MsgType {
		// 消息转发
		case int32(global.RpcMsgTypeTransfer):
			data := msg.Data
			logrus.Infof("收到转发消息：%v\n", data)
		}
	}
}
