package tcp2udp
import (
	"net"
	"log"
	"sync"
	"udptunnel"
	"udpsession"
	"udppacket"
)

const maxBundleClientConns = 0x10
const maxLongConns = 10
const serverAddr = "localhost:9001"

var lock *sync.Mutex
var sessionCount uint32				// 产生sessionId
var ut *udptunnel.UDPTunnel
var idSessionMap map[uint32]*udpsession.Session
func releaseSession(id uint32, flag bool) {
	s, ok := idSessionMap[id]
	if ok {
		lock.Lock()
		defer lock.Unlock()
		s, ok = idSessionMap[id]
		if ok {
			delete(idSessionMap, id)
			p := s.Destroy(flag)
			if p != nil {
				ut.WritePacketToServerProxy(p)
			}
		}
	}
}

/**
 * tunnel call 
 *
 **/
func onData (p *udppacket.Packet) int {
	log.Println("tcp2udp onData")
	p.LogPacket()
	// getsession
	s, ok := idSessionMap[p.SessionId]
	if !ok {
		return -1
	}
	if (p.Length == 0) {
		releaseSession(s.GetId(), false)
	}
	// processNewPacketFromServerProxy
	s.Slock.Lock()
	s.ProcessNewPacketFromServerProxy(p)
	// getNextDataToSend

	for {
		p := s.GetNextRecvDataToSend()
		if p == nil {
			break
		}
		processWrite(s, p.GetPacket())

	}
	s.Slock.Unlock()
	return 0
}

/**
 * client tcp write
 **/
func processWrite (s *udpsession.Session,data []byte) int {
	conn := *s.C
	id := s.GetId()
	index := 0

	for {
		length, err := conn.Write(data[index:])
		if err != nil {
			releaseSession(id, true)
			return -1
		}
		if length != len(data) {
			index += length
		} else {
			break
		}
	}
	return 0
}


/**
 * client tcp read
 **/
func processRead (s *udpsession.Session) {
	log.Println("tcp2udp processRead")
	conn := *s.C
	id := s.GetId()
	for {
		/////////////////////////////////////////////////
		buf := make([]byte, 4096)
		length, err := conn.Read(buf[96:])
		if err != nil {
			log.Println("client read error", err)
			releaseSession(id, true)
			break
		}
		/////////////////////////////////////////////////
		s.ProcessNewDataToServerProxy(buf[:length + 96])

		for {
			p := s.GetNextSendDataToSend()
			if p == nil {
				break
			}
			p.LogPacket()
			rc := ut.WritePacketToServerProxy(p)
			// 检查数据处理结果
			if rc == -1 {
				log.Println("client send err")
				break
			}
		}
	}
	conn.Close()
}


/**
 * client tcp accept
 * one thread call
 **/
func processNewAcceptedConn(conn net.Conn) *udpsession.Session {
	s := udpsession.CreateNewSession(sessionCount)
	s.C = &conn
	idSessionMap[sessionCount] = s
	sessionCount++
	log.Println("sessionCount", sessionCount)
	return s
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
		s := processNewAcceptedConn(conn)
		// load balance, then process conn
		go processRead(s)
	}
}


func Run() {
	sessionCount = 0
	lock = new(sync.Mutex)
	idSessionMap = map[uint32]*udpsession.Session{}
	ut = udptunnel.CreateClientTunnel(onData)
	initListen()
}
