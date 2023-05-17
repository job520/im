package service

import "im/websocket/variables"

func ReceiveGroupMsg(groupID int, msg string) {
	for _, v := range variables.ClientMap {
		if v.GroupSets.Has(groupID) {
			v.DataQueue <- msg
		}
	}
}
