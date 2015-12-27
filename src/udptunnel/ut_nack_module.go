package udptunnel

import (
	"log"
	"sync"
	"udppacket"
)

/*****************************************
 *	sack: 在发包的时候进行判断，是否增加sack
 *  - 发包：添加一个最小的NACK的包，可以设置是否顺便ACK
 *  - 收包：根据包的状态去判断是否ACk， 
 *****************************************/

type UT_NACK struct {
	Index 		   int
	lock           *sync.Mutex 
	Ut             *UDPTunnel
    LOG            *log.Logger
    
    SendMap        map[uint32]*udppacket.Packet
    minSendId      uint32
    maxSendId      uint32
    
    RecvMap        map[uint32]*udppacket.Packet
    minRecvId      uint32
    maxRecvId      uint32 
}

type UT_NACK_CTX struct {
    status         int
} 

func NewNack(LOG *log.Logger) *UT_NACK {
	utNack := new(UT_NACK)
    utNack.Index = -1
    utNack.LOG = LOG
    utNack.SendMap = map[uint32]*udppacket.Packet{}
    utNack.minSendId = 0
    utNack.maxSendId = 0
    utNack.RecvMap = map[uint32]*udppacket.Packet{}
    utNack.minRecvId = 0
    utNack.maxRecvId = 0
	return utNack
}
/********************* interface implements ***********************/
func (utNack *UT_NACK)Debug() string {
	return string(utNack.Index)
}

func (utNack *UT_NACK)InitHandler(ut *UDPTunnel, index int) {
    utNack.Index = index
	utNack.lock = new(sync.Mutex) 
    utNack.Ut = ut
	utNack.LOG.Print("InitHandler")
}

func (utNack *UT_NACK)WriteToServerProxy(p *udppacket.Packet) *udppacket.Packet {
	utNack.LOG.Println("ut_nack_module WriteToServerProxy")
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	utNack.SendMap[p.TunnelId] = p
    p.ChangeOtherTunnelId(utNack.getNackPacketTunnelId())
	return p
}

func (utNack *UT_NACK)WriteToClientProxy(p *udppacket.Packet) *udppacket.Packet {
	utNack.LOG.Println("ut_nack_module WriteToClientProxy")
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	utNack.SendMap[p.TunnelId] = p
    p.ChangeOtherTunnelId(utNack.getNackPacketTunnelId())
	return p
}
 
func (utNack *UT_NACK)ReadFromServerProxy(p *udppacket.Packet) *udppacket.Packet {
	utNack.LOG.Println("ut_nack_module ReadFromServerProxy")
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	nackTunnelId := p.OtherTunnelId
    utNack.ackPackets(nackTunnelId)
    
	return p
}

func (utNack *UT_NACK)ReadFromClientProxy(p *udppacket.Packet) *udppacket.Packet {
	utNack.LOG.Println("ut_nack_module ReadFromClientProxy")
	utNack.lock.Lock()
    defer utNack.lock.Unlock()
	nackTunnelId := p.OtherTunnelId
    utNack.ackPackets(nackTunnelId)
	return p
}

/************** self function ****************/

func (utNack *UT_NACK) getNackPacketTunnelId() uint32 {
    for i := utNack.minRecvId; i < utNack.maxRecvId; i++ {
        if _, ok := utNack.RecvMap[i]; ok {
            delete(utNack.RecvMap, i)
            utNack.minRecvId = i + 1
        } else {
            break
        }
    }
    return utNack.minRecvId
}

func (utNack *UT_NACK) ackPackets(nackTunnelId uint32) {
    for i := utNack.minSendId; i < nackTunnelId; i++ {
        if _, ok := utNack.SendMap[i]; ok {
            delete(utNack.SendMap, i)
            utNack.minSendId = i + 1
        } else {
            utNack.LOG.Print("ackPackets error:missing send packet to ack")
        }
    }
}

func (utNack *UT_NACK) nackPacket(tunnelId uint32) *udppacket.Packet {
    return nil
}
