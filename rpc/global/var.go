package global

import (
	"im/rpc/generate/transfer"
	"sync"
)

type connectMap struct {
	ClientMap map[string]transfer.Transfer_ChatServer
	sync.Mutex
}

var ConnectMap = connectMap{
	ClientMap: make(map[string]transfer.Transfer_ChatServer),
}

const (
	RpcMsgTypeHeartBeat int = iota
	RpcMsgTypeTransfer
)
