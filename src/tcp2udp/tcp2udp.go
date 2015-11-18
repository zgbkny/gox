package tcp2udp
import "net"
import "log"
import "udptunnel"

const maxBundleClientConns = 0x10
const maxLongConns = 10
const serverAddr = "localhost:9001"


var ut *udptunnel.UDPTunnel

/* 处理连接的程序产生对应数据时回调该函数 */
func onData (conn *net.Conn, data []byte) int {
	return processWrite(*conn, data)
}

func processWrite (conn net.Conn, data []byte) int {
	log.Println("processWrite", string(data))
	index := 0
	for {
		length, err := conn.Write(data[index:])
		if err != nil {
			conn.Close()
			return -1
		}
		if length != len(data) {
			index = length
		} else {
			break
		}
	}
	return 0
}
/**
 * 客户端连接的初始
 **/
func processRead (conn net.Conn) {
	log.Println("new connection:", conn.LocalAddr())

	for {
		buf := make([]byte, 4096)
		length, err := conn.Read(buf[96:])
		if err != nil {
			log.Println("client read error", err)
			ut.ProcessCloseConn(conn)
			break
		}
		// 检查数据处理结果
		rc := ut.WriteRawDataToServer(conn, buf[:length + 96], "localhost:90")
		if rc == -1 {
			log.Println("client send err")
			break
		}
	}
	conn.Close()
}


func initListen() {
	// create listener
	listener, err := net.Listen("tcp", "0.0.0.0:9000")
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
		go processRead(conn)
	}
}


func Run() {
	ut = udptunnel.CreateClientTunnel(onData)
	initListen()
}
