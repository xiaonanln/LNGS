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
	j := readRpcMessage(conn)
	fmt.Printf("%v", j)
}

func readRpcMessage(conn net.Conn) int {
	packetBuf := make([]byte, 1024)
	lengthBuf := packetBuf[0:4]
	payloadBuf := packetBuf[4:]

	readAllBytes(conn, lengthBuf)
	length := int(lengthBuf[0]) + int(lengthBuf[1])*256 + int(lengthBuf[2])*256*256 + int(lengthBuf[3])*256*256*256
	fmt.Println("read packet length ", length)
	if length > len(payloadBuf) {
		// error, length too long

	}
	return 1
}

func readAllBytes(conn net.Conn, buff []byte) error {
	if len(buff) <= 0 {
		return nil
	}

	left := len(buff)
	for left > 0 {
		nr, err := conn.Read(buff)
		fmt.Printf("conn.Read %v %v left %s", nr, err, left-nr)
		if err != nil {
			return err
		}
		buff = buff[nr:]
		left = left - nr
	}
	return nil
}
