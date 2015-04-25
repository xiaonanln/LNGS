package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5000")
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept failed: %v", err)
			continue
		}

		serve(conn)
	}

}

func serve(conn net.Conn) {
	fmt.Printf("connection from %s", conn.RemoteAddr().String())

}
