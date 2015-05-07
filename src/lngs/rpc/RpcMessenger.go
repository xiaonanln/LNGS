package lngsrpc

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
)

type Message map[string]interface{}

type MessageEncoder interface {
	EncodeMessage(msg Message) ([]byte, error)
	DecodeMessage(data []byte) (Message, error)
}

type JsonMessageEncoder struct{}
type BsonMessageEncoder struct{}

func (msgEncoder JsonMessageEncoder) EncodeMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

func (msgEncoder JsonMessageEncoder) DecodeMessage(data []byte) (Message, error) {
	msg := make(Message)
	log.Println("decode message %d bytes", len(data))
	err := json.Unmarshal(data, msg)
	return msg, err
}

func (msgEncoder BsonMessageEncoder) EncodeMessage(msg Message) ([]byte, error) {
	return bson.Marshal(msg)
}

func (msgEncoder BsonMessageEncoder) DecodeMessage(data []byte) (Message, error) {
	msg := Message{}
	err := bson.Unmarshal(data, msg)
	return msg, err
}

type RPCMessenger struct {
	conn         net.Conn
	msgEncoder   MessageEncoder
	disconnected bool
	lastError    error
}

func NewRPC(conn net.Conn) *RPCMessenger {
	rpc := RPCMessenger{conn, BsonMessageEncoder{}, false, nil}
	return &rpc
}

func (rpc *RPCMessenger) GetLastError() error {
	return rpc.lastError
}

func (rpc *RPCMessenger) RecvMessage() Message {
	lenBytes := make([]byte, 4)
	if !rpc.recvAll(lenBytes) {
		return nil
	}
	payloadLen := uint(lenBytes[0]) + uint(lenBytes[1])<<8 + uint(lenBytes[2])<<16 + uint(lenBytes[3])<<24
	payload := make([]byte, payloadLen) // allocate enough space for payload
	if !rpc.recvAll(payload) {
		return nil
	}

	msg, err := rpc.msgEncoder.DecodeMessage(payload)
	if err != nil {
		log.Printf("RPCMessenger.RecvMessage: decode error: %v", err)
		return nil
	}
	return msg
}

func (rpc *RPCMessenger) SendMessage(msg Message) {
	bytes, err := rpc.msgEncoder.EncodeMessage(msg)
	if err != nil {
		// failed to encode message, send message failed
		log.Printf("RPCMessenger.SendMessage failed: %v", err)
		return
	}
	payloadLen := uint(len(bytes))
	lenBytes := []byte{
		byte(payloadLen & 0xFF), byte((payloadLen >> 8) & 0xFF), byte((payloadLen >> 16) & 0xFF), byte((payloadLen >> 24) & 0xFF),
	}
	rpc.sendAll(lenBytes)
	rpc.sendAll(bytes)
	log.Printf("rpc send msg %v, payload length %d", msg, payloadLen)
}

func (rpc *RPCMessenger) IsDisconnected() bool {
	return rpc.disconnected
}

func (rpc *RPCMessenger) sendAll(bytes []byte) {
	for len(bytes) > 0 {
		n, err := rpc.conn.Write(bytes)
		if err != nil {
			rpc.onNetworkError(err)
			return
		}
		bytes = bytes[n:] // continue to send the left bytes
	}
}

func (rpc *RPCMessenger) recvAll(buff []byte) bool {
	if len(buff) <= 0 {
		return true
	}

	for len(buff) > 0 {
		nr, err := rpc.conn.Read(buff)
		if err != nil {
			rpc.onNetworkError(err)
			return false
		}
		buff = buff[nr:]
	}
	return true
}

func (rpc *RPCMessenger) onNetworkError(err error) {
	rpc.disconnected = true
	rpc.lastError = err
	log.Printf("disconnected, network error: %v", err)
}

func (rpc *RPCMessenger) Disconnect() {
	rpc.disconnected = true
	rpc.conn.Close()
}
