package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"im/tcp/config"
	"im/tcp/global"
	"im/tcp/logic"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	Address = "localhost:9011"
)

// 处理函数
func process(conn net.Conn) {
	defer conn.Close() // 关闭连接

	// 接收消息
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := buf[:n]
		fmt.Println("message from client:", string(data))
		msg := global.Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		token := msg.Token
		fmt.Println("token:", token)
		// 校验 token 是否合法
		userId, platform, islegal := logic.CheckToken(token)
		if !islegal {
			conn.Close()
			return
		}
		global.ConnectMap.Lock()
		mapKey := fmt.Sprintf("%s:%d", userId, platform)
		global.ConnectMap.ClientMap[mapKey] = conn
		global.ConnectMap.Unlock()
		// 处理消息
		logic.Dispatch(userId, platform, msg) // 消息处理
	}
}

func main() {
	quit := make(chan os.Signal, 1) // 退出信号
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	srv, err := net.Listen("tcp", config.Config.Server.Address)
	if err != nil {
		fmt.Println("Listen() failed, err: ", err)
		return
	}
	go shutdown(quit, srv)
	logic.RegisterWsServer() // 服务注册
	go logic.RpcClient()     // 连接到 rpc 服务器（transfer 服务器）
	fmt.Printf("server running at:%s \n", config.Config.Server.Address)
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
