package rpc

import (
	"encoding/json"
	"net"
)

type Message map[string]interface{}

type RpcMessenger struct {
	conn         net.Conn
	disconnected bool
}

func (*RpcMessenger) RecvMessage() Message {
	return make(Message)
}

func (*RpcMessenger) SendMessage(msg *Message) {

}

func (rm *RpcMessenger) IsDisconnected() bool {
	return rm.disconnected
}
