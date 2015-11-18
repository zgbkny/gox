package udp2tcp
import (
	"log"
	"net"
//	"os"
	"udptunnel"
)

var ut *udptunnel.UDPTunnel

var idConnMap map[uint32]*net.Conn

/**
 * 写数据到服务器
 **/
func onData(conn *net.Conn, data []byte) int {
	log.Println("serverproxy ondata", string(data))
	c := *conn
	_, err := c.Write(data)
	if err != nil {
		return -1
	}
	return 0
}

/**
 * 从服务器读取数据
 **/
func loopRead(c *net.Conn, id uint32) {
	ut.AddNewConnId(*c, id)
	for {
		buf := make([]byte, 4096)
		n, err := (*c).Read(buf[96:])
		if err != nil {
			break
		}
		log.Println("server proxy loopRead", string(buf[96 : 96 + n]))
		ut.WriteRawDataToClient(*c, buf[:96 + n])
	}
}


func Run() {
	idConnMap = map[uint32]*net.Conn{}

	// 启动udp服务器代理，并注册响应的回调函数
	ut = udptunnel.CreateServerTunnel(onData, loopRead)
	ut.StartServer()
}

