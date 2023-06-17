package logic

import (
	"fmt"
	"im/rpc/generate/transfer"
	"io"
	"time"
)

type TransferService struct {
	requests []*transfer.PingResponse
}

func (t TransferService) Ping(conn transfer.Transfer_PingServer) error {
	go func() {
		for {
			request := &transfer.PingResponse{}
			request.Message = "pong!"
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
		fmt.Printf("connector is: %s ,message is: %s .\n", msg.Connector, msg.Message)
	}
}

func (t TransferService) Transfer(conn transfer.Transfer_TransferServer) error {
	go func() {
		for {
			request := &transfer.TransferRequestAndResponse{
				FromConnector: "xx:80",
				ToConnector:   "xx:81",
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
