package main

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"im/rpc/config"
	"im/rpc/generate/transfer"
	"im/rpc/logic"
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
	transferService := logic.TransferService{}
	transfer.RegisterTransferServer(srv, transferService)
	reflection.Register(srv) // 注册到grpcurl
	logrus.Info("Listen on " + config.Config.Server.Address)
	srv.Serve(listen)
}

func shutdown(quit chan os.Signal, srv *grpc.Server) {
	<-quit
	logrus.Info("server shutdown...")
	srv.Stop()
	os.Exit(0)
}
