package main

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"im/rpc/config"
	"im/rpc/generate/transfer"
	"io"
	"time"
)

func main() {
	srv, err := grpc.Dial(config.Config.Server.Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer srv.Close()
	mdClient := metadata.Pairs(
		"fromConnector", ":80",
	)
	ctx := metadata.NewOutgoingContext(context.Background(), mdClient)
	c := transfer.NewTransferClient(srv)
	conn, err := c.Chat(ctx)
	if err != nil {
		logrus.Error("create stream error:", err)
	}
	go func() {
		for {
			if err := conn.Send(&transfer.ChatRequestAndResponse{
				FromConnector: ":80",
				Message:       "ping!",
			}); err != nil {
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()
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
	}
}
