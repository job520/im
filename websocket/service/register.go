package service

import (
	"context"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func KeepAliveRegister(client *clientv3.Client, ctx context.Context, leaseGrant *clientv3.LeaseGrantResponse, timeout int64) {
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
