package udp2tcp
import (
	"log"
	"net"
//	"os"
	"udptunnel"
	"udpsession"
	"udppacket"
)

var ut *udptunnel.UDPTunnel
var idSessionMap map[uint32]*udpsession.Session


func connectToServer(s *udpsession.Session) bool {
	conn, err := net.Dial("tcp", "localhost:90")
	if err != nil {
		log.Println("connectToServer", err)
		return false
	}
	s.C = &conn
	return true
}

/**
 * tunnel call
 **/
func onData(data []byte) int {
	log.Println("udp2tcp onData")
	// getsession
	p := udppacket.GenPacketFromData(data)
	if p == nil {
		return -1
	}
	s, ok := idSessionMap[p.SessionId]
	if !ok {
		log.Println("nil session", p.SessionId)
		s = udpsession.CreateNewSession(p.SessionId)
		ret := connectToServer(s)
		if !ret {
			log.Println("connect to proxy error")
			s.Destroy()
			return 1
		}
		go processRead(s)
	}

	s.ProcessNewPacketFromClientProxy(p)
	for {
		p := s.GetNextRecvDataToSend()
		if p == nil {
			log.Println("send data nil")
			break
		}
		log.Println("udp2tcp processWrite", string(p.GetPacket()))
		processWrite(s, p.GetPacket())

	}
	return 0
}
/**
 * server write
 **/
func processWrite(s *udpsession.Session, data []byte) {
	conn := *s.C
	index := 0
	for {
		length, err := conn.Write(data[index:])
		if err != nil {
			conn.Close()
			return
		}
		if length != len(data) {
			index = length
		} else {
			break
		}
	}
	return
}

/**
 * server read
 **/
func processRead(s *udpsession.Session) {
	conn := *s.C
	log.Println("udp2tcp processRead")
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf[96:])
		if err != nil {
			break
		}

		s.ProcessNewDataToClientProxy(buf[:96 + n])
		log.Println("server proxy loopRead len", n)
		for {
			p := s.GetNextSendDataToSend()
			if p == nil {
				break
			}
			// tunnel write
			ut.WritePacketToClientProxy(p.GetPacket())
		}
	}
}


func Run() {
	idSessionMap = map[uint32]*udpsession.Session{}

	// 启动udp服务器代理，并注册响应的回调函数
	ut = udptunnel.CreateServerTunnel(onData)
	ut.StartServer()
}

