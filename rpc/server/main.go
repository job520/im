package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"im/rpc/proto/hello"
	"io"
	"net"
	"time"
)

const (
	Address = "127.0.0.1:50052"
)

type helloService struct {
	requests []*hello.HelloResponse
}

var HelloService = helloService{}

func (h helloService) SayHello(stream hello.Hello_SayHelloServer) error {
	go func() {
		for {
			request := &hello.HelloResponse{}
			request.Message = "hello from server!"
			if err := stream.Send(request); err != nil {
				fmt.Println("send err:", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("receive done!")
			return nil
		}
		if err != nil {
			fmt.Println("receive error:", err)
			return err
		}
		fmt.Println("msg from client:", msg.Name)
	}
}
func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hello.RegisterHelloServer(s, HelloService)
	reflection.Register(s) // 注册到grpcurl
	fmt.Println("Listen on " + Address)
	s.Serve(listen)
}
