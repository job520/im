package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var msgChan = make(chan string)

func wsClient() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	// todo: 填写 token
	token := ""

	socketUrl := "ws://localhost:8091" + "/chat"
	header := http.Header{}
	header.Add("token", token)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, header)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 发送消息
	go func() {
		for {
			select {
			case msg := <-msgChan:
				err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					fmt.Println("Error during writing to websocket:", err)
				}
			case <-quit:
				fmt.Println("control + c pressed!")
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					fmt.Println("Error during closing websocket:", err)
					os.Exit(0)
				}
				os.Exit(0)
			}
		}
	}()

	// 接收消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error in receive:", err)
			return
		}
		fmt.Println("msg from server:", string(msg))
	}
}

func main() {
	go wsClient()
	// 将要发送到服务端的消息传递到消息管道
	for i := 0; i < 100; i++ {
		msg := `
				{
					"id": 2,
					"cmd": 1,
					"destID": "647acd26c257bfdac6a1c494",
					"msg": "hello i am lee",
					"msgType": 1,
					"ackMsgID": 1
				}
`
		msgChan <- msg
		time.Sleep(3 * time.Second)
	}
	select {}
}
