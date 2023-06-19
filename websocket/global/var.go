package global

import (
	"github.com/gorilla/websocket"
	"im/websocket/generate/transfer"
	"sync"
)

type connectMap struct {
	ClientMap map[string]*websocket.Conn
	sync.Mutex
}

var ConnectMap = connectMap{
	ClientMap: make(map[string]*websocket.Conn),
}

type Message struct {
	FromId string `json:"fromId,omitempty" form:"fromId"` // 发送消息用户ID
	ToId   string `json:"ToId,omitempty" form:"ToId"`     // 接收消息用户ID
	Msg    string `json:"msg,omitempty" form:"msg"`       // 消息内容
}

var RpcMsgChan = make(chan *transfer.ChatRequestAndResponse)
