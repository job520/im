package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"im/http/driver"
	"im/http/global"
	"log"
	"math/rand"
	"time"
)

func GetGateWay(ctx *gin.Context) (string, error) {
	client, err := driver.NewEtcdClient()
	if err != nil {
		logrus.Errorf("etcd init error,%v\n", err)
	}
	key := global.EtcdWsDir
	getResponse, err := client.Get(ctx, key, clientv3.WithPrefix()) // 以前缀获取
	if err != nil {
		log.Printf("etcd GET error,%v\n", err)
		return "", err
	}
	serverLen := len(getResponse.Kvs)
	randIndex := 0
	rand.Seed(time.Now().Unix()) // 播种
	if serverLen > 0 {
		randIndex = rand.Intn(serverLen)
	}
	for index, kv := range getResponse.Kvs {
		if index == randIndex {
			return string(kv.Value), nil
		}
	}
	return "", errors.New("未找到合适的服务器")
}
