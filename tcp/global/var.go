package global

import (
	"im/tcp/generate/transfer"
	"net"
	"sync"
)

type connectMap struct {
	ClientMap map[string]net.Conn
	sync.Mutex
}

var ConnectMap = connectMap{
	ClientMap: make(map[string]net.Conn),
}

type Message struct {
	Token  string `json:"token" form:"token"`             // jwt token
	FromId string `json:"fromId,omitempty" form:"fromId"` // 发送消息用户ID
	ToId   string `json:"ToId,omitempty" form:"ToId"`     // 接收消息用户ID
	Msg    string `json:"msg,omitempty" form:"msg"`       // 消息内容
}

var RpcMsgChan = make(chan *transfer.ChatRequestAndResponse)
