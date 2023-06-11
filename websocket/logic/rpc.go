package logic

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"im/websocket/config"
	"im/websocket/proto/hello"
	"io"
	"time"
)

func RpcClient() {
	address := config.Config.RpcServer.Address
	logrus.Info("rpc server address:", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logrus.Error(err)
	}
	defer conn.Close()
	c := hello.NewHelloClient(conn)
	ctx := context.Background()
	stream, err := c.SayHello(ctx)
	if err != nil {
		fmt.Println("create stream error:", err)
	}
	go func() {
		for {
			if err := stream.Send(&hello.HelloRequest{
				Name: "hello from client!",
			}); err != nil {
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("receive done!")
			break
		}
		if err != nil {
			fmt.Println("receive error:", err)
			break
		}
		if msg != nil {
			fmt.Println("msg from server:", msg.Message)
		}
	}
}
