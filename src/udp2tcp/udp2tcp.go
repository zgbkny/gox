package udp2tcp
import (
	"log"
	"net"
	"os"
	"sync"
	"udptunnel"
	"udpsession"
	"udppacket"
)

var LOG *log.Logger
var ut *udptunnel.UDPTunnel
var idSessionMap map[uint32]*udpsession.Session
var lock *sync.Mutex


func getSession(id uint32) *udpsession.Session {
	s, ok := idSessionMap[id]
	if !ok {
		lock.Lock()
		s, ok = idSessionMap[id]
		if !ok {
			s = udpsession.CreateNewSession(id, LOG)
            s.ModulesCount = ut.ModulesCount
			idSessionMap[id] = s
			ok := connectToServer(s)
			if !ok {
				delete (idSessionMap, id)
				s.Destroy(false)
				lock.Unlock()
				return nil
			}
			lock.Unlock()
			go processRead(s)
		}
	}
	return s
}

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
				ut.WritePacketToClientProxy(p)
			}
		}
	}
}

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
func onData(p *udppacket.Packet) int {
	LOG.Println("udp2tcp onData")
	p.LogPacket()
	// close packet
	if p.Length == 0 {
		releaseSession(p.SessionId, false)
		return 1
	}
	s := getSession(p.SessionId)
	if s == nil {
		return -1
	}
	//LOG.Println("udp2tcp check")
	s.Slock.Lock()
	s.ProcessNewPacketFromClientProxy(p)
	for {
		p := s.GetNextRecvDataToSend()
		if p == nil {
			LOG.Println("send data nil")
			break
		}
		processWrite(s, p.GetPacket())

	}
	s.Slock.Unlock()
	return 0
}
/**
 * server write
 **/
func processWrite(s *udpsession.Session, data []byte) {
	conn := *s.C
	index := 0
	id := s.GetId()
	for {
		length, err := conn.Write(data[index:])
		if err != nil {
			releaseSession(id, true)
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
	LOG.Println("udp2tcp processRead")
	id := s.GetId()
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf[96:])
		if err != nil {
			LOG.Println("server read ", err)
			releaseSession(id, true)
			break
		}

		s.ProcessNewDataToClientProxy(buf[:96 + n])
		for {
			p := s.GetNextSendDataToSend()
			if p == nil {
				break
			}
			// tunnel write
			p.LogPacket()
			ut.WritePacketToClientProxy(p)
		}
	}
}


func Run() {
    
    fileName := "udp2tcp_debug.log"
    logFile,err  := os.Create(fileName)
    defer logFile.Close()
    if err != nil {
        LOG.Fatalln("open file error !")
    }
    LOG = log.New(logFile,"[Debug]",log.Llongfile)
    
	idSessionMap = map[uint32]*udpsession.Session{}
	lock = new(sync.Mutex)
	// 启动udp服务器代理，并注册响应的回调函数
	ut = udptunnel.CreateServerTunnel(onData, LOG)
	ut.Handlers = make([]udptunnel.TunnelHandler, 0)
    utAddId := udptunnel.NewUtAddId(LOG)
    utAddId.LOG = LOG
    ut.Handlers = append(ut.Handlers, utAddId)
	nackModule := udptunnel.NewNack(LOG) 
	ut.Handlers = append(ut.Handlers, nackModule)
    pacingModule := udptunnel.NewPacing(LOG)
	ut.Handlers = append(ut.Handlers, pacingModule)
    utWriter := udptunnel.NewUtWriter(LOG)
    ut.Handlers = append(ut.Handlers, utWriter)
	ut.InitHandlers()
    
	ut.StartServer()
}

