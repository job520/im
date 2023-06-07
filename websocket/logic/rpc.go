package logic

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"im/websocket/config"
	"im/websocket/proto/hello"
	"io"
	"time"
)

func RpcClient() {
	conn, err := grpc.Dial(config.Config.RpcServer.Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
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
		}
		fmt.Println("msg from server:", msg.Message)
	}
}
