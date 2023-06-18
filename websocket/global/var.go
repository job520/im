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

// 定义命令行格式
const (
	HeartMsg int = iota
	SingleMsg
	GroupMsg
)

type Message struct {
	ID       int    `json:"id,omitempty" form:"id"`             // 消息ID
	CMD      int    `json:"cmd,omitempty" form:"cmd"`           // 消息类型（单聊/群聊/心跳）
	FromID   string `json:"fromID,omitempty" form:"fromID"`     // 发送消息用户ID
	DestID   string `json:"destID,omitempty" form:"destID"`     // 接收消息用户ID
	Msg      string `json:"msg,omitempty" form:"msg"`           // 消息内容
	MsgType  int    `json:"msgType,omitempty" form:"msgType"`   // 消息自定义类型（文本消息/图片消息/语音消息...）
	AckMsgID int    `json:"ackMsgID,omitempty" form:"ackMsgID"` // 回复消息ID
}

var RpcMsgChan = make(chan *transfer.ChatRequestAndResponse)

const (
	RpcMsgTypeHeartBeat int = iota
	RpcMsgTypeTransfer
)
