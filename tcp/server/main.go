package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serverMsgId uint32 = iota
	clientMsgId
)

type serverRouter struct {
	znet.BaseRouter
}

func (r *serverRouter) Handle(request ziface.IRequest) {
	// 接收消息
	data := request.GetData()
	fmt.Println("message from client:", string(data))
}

func main() {
	quit := make(chan os.Signal, 1) // 退出信号
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	//srv := znet.NewUserConfServer(&zconf.Config{
	//    Host:              "127.0.0.1",
	//    TCPPort:           8999,
	//    LogIsolationLevel: 2,
	//})
	srv := znet.NewServer()
	srv.SetOnConnStart(func(conn ziface.IConnection) {
		// 发送消息
		go func() {
			for {
				err := conn.SendMsg(clientMsgId, []byte("hello from server!"))
				if err != nil {
					fmt.Println(err)
					break
				}
				time.Sleep(3 * time.Second)
			}
		}()
	})
	srv.AddRouter(serverMsgId, &serverRouter{})
	go shutdown(quit, srv)
	srv.Serve()
}

func shutdown(quit chan os.Signal, srv ziface.IServer) {
	<-quit
	fmt.Println("server shutdown...")
	srv.Stop()
	os.Exit(0)
}
