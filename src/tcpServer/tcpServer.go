package tcpServer
import (
)
/*

func initListen(l * tunnel.TunnelListener) {
	
	// create listener
	listener, err := net.Listen("tcp", "0.0.0.0:9001")
	
	if err != nil {
		return
	}
	
	// listen and accept connections from clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		// create a tunnel for the conn
		go tunnel.ProcessTunnel(conn)
	}
}

func Run() {
	l := new(tunnel.TunnelListener)
	
	// 初始化监听，for accept tunnel
	initListen(l)
}*/
