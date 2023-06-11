package logic

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/rpc/config"
	"im/rpc/driver"
	"im/rpc/global"
	"im/rpc/proto/hello"
	"io"
	"time"
)

type HelloService struct {
	requests []*hello.HelloResponse
}

func (h HelloService) SayHello(stream hello.Hello_SayHelloServer) error {
	go func() {
		for {
			request := &hello.HelloResponse{}
			request.Message = "hello from server!"
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
		fmt.Println("msg from client:", msg.Name)
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
