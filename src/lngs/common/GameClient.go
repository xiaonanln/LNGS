package lngscommon

import (
	. "lngs/rpc"
	"log"
	"net"
)

type GameClient struct {
	rpc   *RPCMessenger
	owner *Entity
}

func NewGameClient(conn net.Conn) *GameClient {
	client := GameClient{NewRPC(conn), nil}
	return &client
}

func (self *GameClient) RecvMessage() Message {
	return self.rpc.RecvMessage()
}

func (self *GameClient) IsDisconnected() bool {
	return self.rpc.IsDisconnected()
}

func (self *GameClient) Disconnect() {
	self.rpc.Disconnect()
}

func (self *GameClient) OnReceiveMessage(msg Message) {
	if self.owner == nil {
		// owner is nil, message droped
		log.Printf("message droped: %v", msg)
		return
	}

	self.owner.OnReceiveMessage(msg)
}
