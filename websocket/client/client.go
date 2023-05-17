package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func receive(conn *websocket.Conn) {
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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	socketUrl := "ws://localhost:8080" + "/chat?id=1&token=xxx"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	go receive(conn)

	for {
		select {
		case <-time.After(3 * time.Second):
			msg := `
				{
					"id": 2,
					"cmd": 1,
					"fromID": 1,
					"destID": 1,
					"msg": "hello from client!",
					"msgType": 1,
					"ackMsgID": 1
				}
`
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Error during writing to websocket:", err)
				return
			}
		case <-quit:
			fmt.Println("control + c pressed!")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("Error during closing websocket:", err)
				return
			}
			return
		}
	}
}
