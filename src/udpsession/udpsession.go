package udpsession

import (
	"net"
//	"log"
	"container/list"
	"utils"
	"sync"
	"udppacket"
)

const SESS_INIT = 0
const SESS_NORMAL = 1
/** 单独一个会话,包含会话的所有信息以及当前的状态 **/
type Session struct {
	id				uint32
	C				*net.Conn
	dst				string

	status			int
	onDataF			func(*net.Conn, []byte) int
	Slock			*sync.Mutex

	sendPackets		[]*udppacket.Packet						// 发送列表

	idPacketRecvMap	map[uint32]*udppacket.Packet			// udptunnel接收
	maxPacketRecvId	uint32								// 从udptunnel recv
	minPacketRecvId	uint32								// 从udptunnel recv

	sendList		*list.List							// 到udptunnel发送
	recvList		*list.List							// 从udptunnel接收
	// 统计
	count			uint32
}

func CreateNewSession(id uint32) *Session {
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
	return s
}


/**
 * destroy session
 * -- return empty packet to ack remote proxy
 */
func (s *Session) Destroy() *udppacket.Packet {
	p := udppacket.CreateNewPacket(0)
	return p
}

func (s *Session) ProcessNewDataToServerProxy(rawData []byte) {
	//log.Println("session ProcessNewDataToServerProxy")
	p := udppacket.CreateNewPacket(s.count, rawData, "")
	p.SessionId = s.id
	s.sendPackets = append(s.sendPackets, p)
	s.count++
	header := s.genHeader(p)
	p.RawDataAddHeader(header)
	s.sendList.PushBack(p)
}

func (s *Session) ProcessNewDataToClientProxy(rawData []byte) {
	//log.Println("session ProcessNewDataToClientProxy")
	p := udppacket.CreateNewPacket(s.count, rawData, "")
	s.sendPackets = append(s.sendPackets, p)
	s.count++
	header := s.genHeader(p)
	p.RawDataAddHeader(header)
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

func (s *Session) genHeader(p *udppacket.Packet) []byte{
	//log.Println("udptunnel genHeader")
	header := make([]byte, 96)
	count := 0
	//log.Println("genHeader datalen", len(p.RawData) - 96)
	dataLenBytes := utils.Int16ToBytes(len(p.RawData) - 96)
	copy(header[count:(count + 2)], dataLenBytes)
	count += 2
	sessionIdBytes := utils.Int32ToBytes(p.SessionId)
	copy(header[count:(count + 4)], sessionIdBytes)
	count += 4
	packetIdBytes := utils.Int32ToBytes(p.Id)
	copy(header[count:(count + 4)], packetIdBytes)
	count += 4
	// 传输层协议类型
	header[count] = p.ProtoType
	count++
	// 包类型
	header[count] = p.PacketType
	count++
	if s.status == SESS_INIT {
		header[count] = byte(len(p.Dst))
		count++
		if len(p.Dst) > 0 {
			dstBytes := []byte(p.Dst)
			copy(header[count:(count + len(dstBytes))], dstBytes)
			count += len(dstBytes)
		}
	} else if s.status == SESS_NORMAL {
		header[count] = byte(SESS_NORMAL)
		count++
	}
	return header[:count]
}
