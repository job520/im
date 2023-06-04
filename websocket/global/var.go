package global

import (
	"github.com/gorilla/websocket"
	"im/websocket/utils"
	"sync"
)

var RwLocker sync.RWMutex

// userid 和 Node 映射关系表
var ClientMap = make(map[string]*Node, 0)

type Node struct {
	Conn *websocket.Conn
	// 并行转串行,
	DataQueue chan string
	GroupSets utils.Set
}

// 定义命令行格式
const (
	HeartMsg int = iota
	SingleMsg
	GroupMsg
)

type Message struct {
	ID       int    `json:"id,omitempty" form:"id"`             // 消息ID
	CMD      int    `json:"cmd,omitempty" form:"cmd"`           // 消息类型（单聊/群聊...）
	FromID   string `json:"fromID,omitempty" form:"fromID"`     // 发送消息用户ID
	DestID   string `json:"destID,omitempty" form:"destID"`     // 接收消息用户ID
	Msg      string `json:"msg,omitempty" form:"msg"`           // 消息内容
	MsgType  int    `json:"msgType,omitempty" form:"msgType"`   // 消息自定义类型（文本消息/图片消息/语音消息...）
	AckMsgID int    `json:"ackMsgID,omitempty" form:"ackMsgID"` // 回复消息ID
}
