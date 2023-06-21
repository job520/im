package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"time"
)

const (
	serverMsgId uint32 = iota
	clientMsgId
)

var msgChan = make(chan string)

type clientRouter struct {
	znet.BaseRouter
}

func (r *clientRouter) Handle(request ziface.IRequest) {
	// 接收消息
	data := request.GetData()
	fmt.Println("message from server:", string(data))
}

func tcpClient() {
	// 创建 client客户端
	client := znet.NewClient("127.0.0.1", 8999)

	client.SetOnConnStart(func(conn ziface.IConnection) {
		// 发送消息
		go func() {
			for {
				select {
				case msg := <-msgChan:
					err := conn.SendMsg(serverMsgId, []byte(msg))
					if err != nil {
						fmt.Println(err)
						break
					}
				}
			}
		}()
	})
	client.AddRouter(clientMsgId, &clientRouter{})
	// 启动客户端
	client.Start()
}

func main() {
	go tcpClient()
	// 将要发送到服务端的消息传递到消息管道
	for i := 0; i < 3; i++ {
		msgChan <- "hello from client!"
		time.Sleep(3 * time.Second)
	}
	select {}
}
