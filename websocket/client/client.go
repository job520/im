package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func httpPost(url string, params url.Values) (string, error) {
	resp, err := http.PostForm(url, params)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	} else {
		if resp.StatusCode == 200 {
			var bodyReader io.ReadCloser = resp.Body
			body, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				return "", err
			} else {
				return string(body), nil
			}
		} else {
			return "", fmt.Errorf("服务器异常")
		}
	}
}

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
	go receive(conn)

	for {
		select {
		case <-time.After(3 * time.Second):
			msg := `
				{
					"id": 2,
					"cmd": 1,
					"fromID": "user-1",
					"destID": "user-2",
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
