package logic

import (
	"fmt"
	"im/rpc/generate/transfer"
	"io"
	"time"
)

type TransferService struct {
	requests []*transfer.ChatRequestAndResponse
}

func (t TransferService) Chat(conn transfer.Transfer_ChatServer) error {
	go func() {
		for {
			request := &transfer.ChatRequestAndResponse{
				FromConnector: "xx:80",
				ToConnector:   "xx:81",
				MsgType:       1,
				Message:       "pong!",
				Data: &transfer.Data{
					Id:         0,
					Cmd:        0,
					FromId:     "",
					DestId:     "",
					Msg:        "",
					MsgType:    0,
					AckMsgType: 0,
				},
			}
			if err := conn.Send(request); err != nil {
				fmt.Println("send err:", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	for {
		msg, err := conn.Recv()
		if err == io.EOF {
			fmt.Println("receive done!")
			return nil
		}
		if err != nil {
			fmt.Println("receive error:", err)
			return err
		}
		fmt.Printf("msg received: %v .\n", msg)
	}
}
