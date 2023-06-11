package logic

import (
	"context"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/websocket/config"
	"im/websocket/driver"
	"im/websocket/global"
	"im/websocket/service"
)

func RegisterWsServer() {
	client, err := driver.NewEtcdClient()
	if err != nil {
		logrus.Errorf("etcd init error,%v\n", err)
	}
	key := global.EtcdWsDir + config.Config.Server.Address
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
	go service.KeepAliveRegister(client, ctx, leaseGrant, timeout) // 续租
}
