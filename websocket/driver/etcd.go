package driver

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/websocket/config"
	"time"
)

func NewEtcdClient(collectionName string) (*clientv3.Client, error) {
	host := config.Config.Etcd.Host
	port := config.Config.Etcd.Port
	username := config.Config.Etcd.Username
	password := config.Config.Etcd.Password
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%d", host, port)},
		DialTimeout: 5 * time.Second,
		Username:    username,
		Password:    password,
	})
	return client, err
}
