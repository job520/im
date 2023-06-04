package service

import "im/websocket/global"

func ReceiveGroupMsg(groupID int, msg string) {
	for _, v := range global.ClientMap {
		if v.GroupSets.Has(groupID) {
			v.DataQueue <- msg
		}
	}
}
