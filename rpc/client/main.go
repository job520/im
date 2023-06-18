package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"im/rpc/config"
	"im/rpc/generate/transfer"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var msgChan = make(chan *transfer.ChatRequestAndResponse)

func rpcClient() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

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

	// 发送消息
	go func() {
		for {
			select {
			case msg := <-msgChan:
				if err := conn.Send(msg); err != nil {
					return
				}
			case <-quit:
				fmt.Println("control + c pressed!")
				err := conn.CloseSend()
				if err != nil {
					fmt.Println("close error:", err)
					os.Exit(0)
				}
				os.Exit(0)
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
	}
}

func main() {
	go rpcClient()
	for i := 0; i < 3; i++ {
		msg := &transfer.ChatRequestAndResponse{
			FromConnector: ":80",
			Message:       "ping!",
		}
		msgChan <- msg
		time.Sleep(3 * time.Second)
	}
	select {}
}
