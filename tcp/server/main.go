package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	Address = "localhost:9011"
)

// 处理函数
func process(conn net.Conn) {
	defer conn.Close() // 关闭连接

	// 发送消息
	go func() {
		for {
			_, err := conn.Write([]byte("hello from server!"))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()
	// 接收消息
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := buf[:n]
		fmt.Println("message from client:", string(data))
	}
}

func main() {
	quit := make(chan os.Signal, 1) // 退出信号
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	srv, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println("Listen() failed, err: ", err)
		return
	}
	go shutdown(quit, srv)
	fmt.Printf("server running at:%s \n", Address)
	for {
		conn, err := srv.Accept() // 监听客户端的连接请求
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		go process(conn) // 启动一个 goroutine 来处理客户端的连接请求
	}
}

func shutdown(quit chan os.Signal, srv net.Listener) {
	<-quit
	fmt.Println("server shutdown...")
	if err := srv.Close(); err != nil {
		fmt.Println("shutdown error:", err)
	}
	os.Exit(0)
}
