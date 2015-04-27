package main

import (
	"lngs/rpc"
	"log"
	"net"
)

func main() {
	localAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12340")
	if err != nil {
		log.Panicln(err)
	}

	remoteAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5000")
	if err != nil {
		log.Panicln(err)
	}

	conn, err := net.DialTCP("tcp", localAddr, remoteAddr)
	if err != nil {
		log.Panicln(err)
	}
	rpc := rpc.NewRPC(conn, rpc.BsonMessageEncoder{})
	msg := rpc.RecvMessage()
	log.Println(msg)
}
