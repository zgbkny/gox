package tcpClient
import "net"
import "log"
import "tunnel"

const maxBundleClientConns = 0x10
const maxLongConns = 10
const serverAddr = "localhost:9001"

/* 处理连接的程序产生对应数据时回调该函数 */
func onData (conn net.Conn, data []byte) {
	processWrite(conn, data)
}

func processWrite (conn net.Conn, data []byte) {
	log.Println("processWrite", string(data))
	index := 0
	for {
		length, err := conn.Write(data[index:])
		if err != nil {
			tunnel.RemoveSession(conn)
			break
		}
		if length != len(data) {
			index = length
		} else {
			break
		}
	}
}

func processRead (conn net.Conn) {
	log.Println("new connection:", conn.LocalAddr())
	
	for {
		buf := make([]byte, 4096)
		length, err := conn.Read(buf)
		if err != nil {
			log.Println("client read error", err)
			break
		}
		// 检查数据处理结果
		rc := tunnel.ProcessData(onData, conn, buf[:length], "destination")
		if rc == -1 {
			log.Println("client send err")
			break
		}
	}
}


func initListen(l *tunnel.Listener) {
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
	l := new(tunnel.Listener)
	tunnel.Init()
	// init long conns 
	/*for i := 0; i < maxLongConns; i++ {
		addr, err1 := net.ResolveTCPAddr("tcp", serverAddr)
		if err1 != nil {
			log.Fatal(err1)
		}

		conn, err2 := net.DialTCP("tcp", nil, addr)
		if err2 != nil {
			log.Fatal(err2)
		}
		log.Println("i:", i)
		b := newBundle(conn)
		l.bundles[i] = b
		l.id[i] = i
		go processBundle(b)
	}*/
	// 初始化监听，for accept
	initListen(l)
}
