package main

import (
	"net"
	"log"
)
/**
 * client tcp read
 **/
func processRead (c *net.Conn) {
	log.Println("tcp2udp processRead")
	conn := *c
	for {
		/////////////////////////////////////////////////
		buf := make([]byte, 4096)
		length, err := conn.Read(buf[96:])
		log.Println("read over")
		if err != nil {
			log.Println("client read error", err)
			break
		}
		log.Println("check")
		/////////////////////////////////////////////////
		log.Println("read data:", string(buf[96:96 + length]))
		conn.Write(buf[96 : 96 + length])
		break
	}
	log.Println("close")
	conn.Close()
}


/**
 * client tcp accept
 * one thread call
 **/
func processNewAcceptedConn(conn net.Conn) {
	log.Println("processNewAcceptedConn")
}

func initListen() {
	// create listener
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return
	}

	// listen and accept connections from clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		// load balance, then process conn
		go processRead(&conn)
	}
}
func main() {
	log.Println("main")
	initListen() 
}