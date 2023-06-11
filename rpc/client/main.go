package main

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"im/rpc/config"
	"im/rpc/proto/hello"
	"io"
	"time"
)

func main() {
	conn, err := grpc.Dial(config.Config.Server.Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()
	c := hello.NewHelloClient(conn)
	ctx := context.Background()
	stream, err := c.SayHello(ctx)
	if err != nil {
		logrus.Error("create stream error:", err)
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
			logrus.Info("receive done!")
			break
		}
		if err != nil {
			logrus.Error("receive error:", err)
		}
		logrus.Info("msg from server:", msg.Message)
	}
}
