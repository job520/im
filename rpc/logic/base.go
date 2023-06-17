package logic

import (
	"fmt"
	"im/rpc/generate/transfer"
	"im/rpc/global"
	"io"
)

type TransferService struct {
	requests []*transfer.ChatRequestAndResponse
}

func (t TransferService) Chat(conn transfer.Transfer_ChatServer) error {
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
		global.ConnectMap.Lock()
		global.ConnectMap.ClientMap[msg.FromConnector] = conn
		global.ConnectMap.Unlock()
	}
}
