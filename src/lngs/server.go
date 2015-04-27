package main

import (
	"fmt"
	"net"

	"lngs/rpc"
)

func main() {
	serverAddrStr := "0.0.0.0:5000"
	fmt.Println("Resolving TCP address: ", serverAddrStr)
	serverAddr, err := net.ResolveTCPAddr("tcp", serverAddrStr)
	if err != nil {
		panic("failed to resolve server addr")
		return
	}
	fmt.Println(serverAddr)
	listener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		panic("failed to listen")
		return
	}

	fmt.Println("server listening on 	, start accepting connections...", serverAddr)
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("accept failed: %v", err)
		return
	}

	serve(conn)
}

func serve(conn net.Conn) {
	fmt.Printf("connection from %s", conn.RemoteAddr().String())
	rpcer := rpc.NewRPC(conn, rpc.BsonMessageEncoder{})
	rpcer.SendMessage(rpc.Message{"1": 1, "2": 2, "hello": "world"})
}
