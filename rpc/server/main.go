package main

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"im/rpc/config"
	"im/rpc/logic"
	"im/rpc/proto/hello"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1) // 退出信号
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	listen, err := net.Listen("tcp", config.Config.Server.Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	go shutdown(quit, srv)
	helloService := logic.HelloService{}
	hello.RegisterHelloServer(srv, helloService)
	reflection.Register(srv)  // 注册到grpcurl
	logic.RegisterRpcServer() // 服务注册
	logrus.Info("Listen on " + config.Config.Server.Address)
	srv.Serve(listen)
}

func shutdown(quit chan os.Signal, srv *grpc.Server) {
	<-quit
	logrus.Info("server shutdown...")
	srv.Stop()
	os.Exit(0)
}
