package lngs

import (
	. "lngs/db"
	. "lngs/rpc"
	"log"
	"net"
	"runtime/debug"
)

func Run(serverAddrStr string) {

	go serveDB()
	go serveDB()
	go serveDB()

	log.Println("Resolving TCP address: ", serverAddrStr)
	serverAddr, err := net.ResolveTCPAddr("tcp", serverAddrStr)
	if err != nil {
		panic("failed to resolve server addr")
		return
	}
	listener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		panic("failed to listen")
		return
	}

	log.Println("Start accepting connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept failed: %v", err)
			return
		}
		go serveConn(conn)
	}
}

func serveDB() {
	dbcon := ConnectDB()
	defer dbcon.Close()
	dbmanager := NewDbManager(dbcon)
	dbmanager.Loop()
}

func serveConn(conn net.Conn) {
	log.Printf("New connection: %s\n", conn.RemoteAddr().String())
	client := NewGameClient(conn)
	serveGameClient(client)
	log.Printf("Connection %s terminated\n", conn.RemoteAddr().String())
}

func serveGameClient(client *GameClient) {
	// create boot entity for the new client
	var bootEntity *Entity = entityManager.NewBootEntity()
	bootEntity.SetClient(client)

	if bootEntity == nil {
		client.Disconnect()
		return
	}

	log.Printf("create boot: %s", bootEntity)
	// boot := *bootEntity.FieldByName("entity")

	for !client.IsDisconnected() {
		msg := client.RecvMessage()
		if msg != nil {
			serveMessage(client, msg)
		}
	}
}

func recoverFromError() {
	if err := recover(); err != nil {
		log.Println("[W]", err)
		debug.PrintStack()
	}
}

func serveMessage(client *GameClient, msg Message) {
	defer recoverFromError()

	log.Printf("Recv message: %v\n", msg)
	client.OnReceiveMessage(msg)
}
