package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	Address = "localhost:9012"
)

var msgChan = make(chan string)

func tcpClient() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	conn, err := net.Dial("tcp", Address)
	if err != nil {
		panic(err)
	}
	defer conn.Close() // 关闭 TCP 连接

	// 发送消息
	go func() {
		for {
			select {
			case msg := <-msgChan:
				_, err := conn.Write([]byte(msg))
				if err != nil {
					fmt.Println("write error:", err)
				}
			case <-quit:
				fmt.Println("control + c pressed!")
				_, err := conn.Write([]byte("client closed!"))
				if err != nil {
					fmt.Println("send close message error:", err)
					os.Exit(0)
				}
				if err := conn.Close(); err != nil {
					fmt.Println("conn close error:", err)
					os.Exit(0)
				}
				os.Exit(0)
			}
		}
	}()
	// 接收消息
	for {
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("Error in receive:", err)
			return
		}
		msg := buf[:n]
		fmt.Println("message from server:", string(msg))
	}
}

func main() {
	go tcpClient()
	// 将要发送到服务端的消息传递到消息管道
	for i := 0; i < 100; i++ {
		// todo: 填写 token
		msgChan <- `
				{
					"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODgzMDAyMDEsInBsYXRmb3JtIjoxLCJ1aWQiOiI2NDdhY2QyNmMyNTdiZmRhYzZhMWM0OTQifQ.dDp2e_rjsDCmMk8pnauJALX3qEvJ6ghp21B1LsN_S3s",
					"toId": "6472bc5228bdc3d600d72b0d",
					"msg": "hello i am lee1"
				}
`
		time.Sleep(3 * time.Second)
	}
	select {}
}
