package lngs

import (
	. "lngs/rpc"
	"log"
	"net"
	"lngs/db"
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

func (self *GameClient) sendMessage(msg Message) {
	self.rpc.SendMessage(msg)
}

func (self *GameClient) CreateEntity(entityid string, entitytype string) {
	self.sendMessage(map[string]interface{}{
		"CE": []string{entityid, entitytype},
	})
}

func (self *GameClient) DestroyEntity(entityid string) {
	self.sendMessage(map[string]interface{}{
		"DE": entityid,
	})
}

func (self *GameClient) Call(entityid string, methodname string, args ...interface{}) {
	self.sendMessage(map[string]interface{}{
		"ID":   entityid,
		"M":    methodname,
		"ARGS": args,
	})
}

func (self *GameClient) BecomePlayer(entityid string, data lngsdb.Doc) {
	self.Call(entityid, "BecomePlayer", data)
}
