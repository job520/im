package logic

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/rpc/config"
	"im/rpc/driver"
	"im/rpc/generate/transfer"
	"im/rpc/global"
	"io"
	"time"
)

type TransferService struct {
	requests []*transfer.PingResponse
}

func (t TransferService) Ping(stream transfer.Transfer_PingServer) error {
	go func() {
		for {
			request := &transfer.PingResponse{}
			request.Message = "pong!"
			if err := stream.Send(request); err != nil {
				fmt.Println("send err:", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	for {
		msg, err := stream.Recv()
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

func (t TransferService) Transfer(stream transfer.Transfer_TransferServer) error {
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
			if err := stream.Send(request); err != nil {
				fmt.Println("send err:", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	for {
		msg, err := stream.Recv()
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

func RegisterRpcServer() {
	client, err := driver.NewEtcdClient()
	if err != nil {
		logrus.Errorf("etcd init error,%v\n", err)
	}
	key := global.EtcdRpcDir + config.Config.Server.Address
	value := config.Config.Server.Address
	timeout := int64(10)
	ctx := context.Background()
	// 获取一个租约 有效期为10秒
	leaseGrant, err := client.Grant(ctx, timeout)
	if err != nil {
		logrus.Errorf("grant error %v", err)
	}
	// PUT 租约期限为10秒
	_, err = client.Put(ctx, key, value, clientv3.WithLease(leaseGrant.ID))
	if err != nil {
		logrus.Errorf("put error %v", err)
	}
	go keepAliveRegister(client, ctx, leaseGrant, timeout) // 续租
}

func keepAliveRegister(client *clientv3.Client, ctx context.Context, leaseGrant *clientv3.LeaseGrantResponse, timeout int64) {
	for {
		// 续租
		keepaliveResponse, err := client.KeepAliveOnce(ctx, leaseGrant.ID)
		if err != nil {
			logrus.Errorf("KeepAlive error %v", err)
			return
		}
		logrus.Info(keepaliveResponse.TTL)
		time.Sleep(time.Duration(timeout/2) * time.Second)
	}
}
