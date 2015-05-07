package main

import (
	"lngs/rpc"
	"log"
	"net"
)

func main() {
	remoteAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:7000")
	if err != nil {
		log.Panicln(err)
	}

	conn, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		log.Panicln(err)
	}
	rpc := lngsrpc.NewRPC(conn)
	for i := 0; i < 10; i++ {
		rpc.SendMessage(map[string]interface{}{"ID": "123", "M": "Test", "ARGS": []int{1, 2, 3}})
	}
}
