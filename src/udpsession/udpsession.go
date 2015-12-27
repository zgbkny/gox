package udpsession

import (
	"net"
	"log"
	"container/list"
//	"utils"
	"sync"
	"udppacket"
)

const SESS_INIT = 0
const SESS_NORMAL = 1
const SESS_CLOSE = 2	// 通告对端
const SESS_RELEASE = 3	// 只释放资源
/** 单独一个会话,包含会话的所有信息以及当前的状态 **/
type Session struct {
    LOG             *log.Logger
	id				uint32
	C				*net.Conn
	dst				string

	status			int
	onDataF			func(*net.Conn, []byte) int
	Slock			*sync.Mutex

	idPacketRecvMap	map[uint32]*udppacket.Packet			// udptunnel接收, order packet
	maxPacketRecvId	uint32								// 从udptunnel recv
	minPacketRecvId	uint32								// 从udptunnel recv

	sendList		*list.List							// 到udptunnel发送
	recvList		*list.List							// 从udptunnel接收
	// 统计
	count			uint32
    ModulesCount    int
}

func CreateNewSession(id uint32, LOG *log.Logger) *Session {
	s := new(Session)
	s.id = id
	s.count = 0
	s.minPacketRecvId = 0
	s.maxPacketRecvId = 0
	s.C = nil
	s.Slock = new(sync.Mutex)
	s.idPacketRecvMap = map[uint32]*udppacket.Packet{}
	s.sendList = list.New()
	s.recvList = list.New()
    s.LOG = LOG
	return s
}
func (s *Session) GetId() uint32 {
	return s.id
}

/**
 * destroy session
 * -- return empty packet to ack remote proxy
 */
func (s *Session) Destroy(flag bool) *udppacket.Packet {
	if s.C != nil {
		conn := *s.C
		conn.Close()
	}
	if flag {
		data := make([]byte, 96)
		p := udppacket.CreateNewPacket(0, data, "", s.ModulesCount, s.LOG)
		p.SessionId = s.id
		p.RawDataAddHeader()
		return p
	} else {
		return nil
	}
}

func (s *Session) ProcessNewDataToServerProxy(rawData []byte) {
	//log.Println("session ProcessNewDataToServerProxy")
	p := udppacket.CreateNewPacket(s.count, rawData, "", s.ModulesCount, s.LOG)
	p.SessionId = s.id
	p.OtherLen = 0
	s.count++
	p.RawDataAddHeader()
	s.sendList.PushBack(p)
}

func (s *Session) ProcessNewDataToClientProxy(rawData []byte) {
	//log.Println("session ProcessNewDataToClientProxy")
	p := udppacket.CreateNewPacket(s.count, rawData, "", s.ModulesCount, s.LOG)
	p.SessionId = s.id
	p.OtherLen = 0
	p.Dst = "" 
	s.count++
	p.RawDataAddHeader()
	s.sendList.PushBack(p)
}
/**
 * 处理从服务端网关发回的数据包
 **/
func (s *Session) ProcessNewPacketFromServerProxy(p *udppacket.Packet) {
	//log.Println("session ProcessNewPacketFromServerProxy")
	if s.maxPacketRecvId < p.Id {
		s.maxPacketRecvId = p.Id
	}
	s.idPacketRecvMap[p.Id] = p
	for {
		p, ok := s.idPacketRecvMap[s.minPacketRecvId]
		if !ok {
			break
		}
		s.recvList.PushBack(p)
		delete(s.idPacketRecvMap, s.minPacketRecvId)
		s.minPacketRecvId++
	}
}

/**
 * 会话处理新的数据包
 * 并发
 **/
func (s *Session) ProcessNewPacketFromClientProxy(p *udppacket.Packet) {
	//log.Println("session ProcessNewPacketFromClientProxy")
	if s.maxPacketRecvId < p.Id {
		s.maxPacketRecvId = p.Id
	}
	s.idPacketRecvMap[p.Id] = p
	//log.Println("p.id", p.Id, "s.minPacketRecvId", s.minPacketRecvId)
	for {
		p, ok := s.idPacketRecvMap[s.minPacketRecvId]
		//log.Println("ok", ok)
		if !ok {
			break
		}
		s.recvList.PushBack(p)
		delete(s.idPacketRecvMap, s.minPacketRecvId)
		s.minPacketRecvId++
	}
}


func (s *Session) GetNextSendDataToSend() *udppacket.Packet {
	e := s.sendList.Front()
	if e != nil {
		if data, ok := e.Value.(*udppacket.Packet); ok {
			s.sendList.Remove(e)
			return data
		}
	}
	return nil
}

func (s *Session) GetNextRecvDataToSend() *udppacket.Packet {
	e := s.recvList.Front()
	if e != nil {
		if data, ok := e.Value.(*udppacket.Packet); ok {
			s.recvList.Remove(e)
			return data
		}
	}
	return nil
}
