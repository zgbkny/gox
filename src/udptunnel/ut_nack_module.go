package udptunnel

import (
	"log"
	"sync"
    "time"
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
    
    SendMap        map[uint64]*PACKET_WRAPPER
    minSendId      uint64
    maxSendId      uint64
    
    RecvMap        map[uint64]*PACKET_WRAPPER
    minRecvId      uint64
    maxRecvId      uint64 
    
    rtrDuration    time.Duration
    ackDuration    time.Duration
    now            *time.Time
}

type CTX_UT_NACK struct {
    status         int 
    rtrTime        *time.Time
} 
const CTX_STATUS_NEW = 1
const CTX_STATUS_RTR = 2      // 重传

/*******************************
 * packet在nack模块中的包装类
 *******************************/
type PACKET_WRAPPER struct {
    status         int
    nackTime       time.Time
    p              *udppacket.Packet // 发送的包
}
const WRAPPER_STATUS_NRECVED = 0x1
const WRAPPER_STATUS_NACKED  = 0x2
const WRAPPER_STATUS_RECVED  = 0x3  
const WRAPPER_STATUS_SENDED  = 0x4
const WRAPPER_STATUS_RTR     = 0x5
const WRAPPER_STATUS_ACKED   = 0x6


const STATUS_NACK  = 0x1
const STATUS_ACK   = 0x2


func NewNack(LOG *log.Logger) *UT_NACK {
	utNack := new(UT_NACK)
    utNack.Index = -1
    utNack.LOG = LOG
    utNack.SendMap = map[uint64]*PACKET_WRAPPER{}
    utNack.minSendId = 0
    utNack.maxSendId = 0
    utNack.RecvMap = map[uint64]*PACKET_WRAPPER{}
    utNack.minRecvId = 0
    utNack.maxRecvId = 0
    tmpTime := time.Now()
    utNack.now = &tmpTime
    utNack.rtrDuration = time.Duration(500000000) // 500ms
    utNack.ackDuration = time.Duration(500000000) // 500ms
	return utNack
}
/********************* interface implements ***********************/
func (utNack *UT_NACK)Debug() string {
	return string(utNack.Index)
}

func (utNack *UT_NACK) InitHandler(ut *UDPTunnel, index int) {
    utNack.Index = index
	utNack.lock = new(sync.Mutex) 
    utNack.Ut = ut
	utNack.LOG.Print("InitHandler")
}



func (utNack *UT_NACK) WriteToServerProxy(p *udppacket.Packet) *udppacket.Packet {
	utNack.LOG.Println("ut_nack_module WriteToServerProxy SendMap size:", len(utNack.SendMap), utNack.minSendId, utNack.maxSendId, "recvMap size:", len(utNack.RecvMap), utNack.minRecvId, utNack.maxRecvId)
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
    utNack.processSendPacket(p)
	return p
}

func (utNack *UT_NACK) WriteToClientProxy(p *udppacket.Packet) *udppacket.Packet {
	utNack.LOG.Println("ut_nack_module WriteToClientProxy SendMap size:", len(utNack.SendMap),  utNack.minSendId, utNack.maxSendId, "RecvMap size:", len(utNack.RecvMap), utNack.minRecvId, utNack.maxRecvId)
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	utNack.processSendPacket(p)
	return p
}

func (utNack *UT_NACK) ReadFromServerProxy(p *udppacket.Packet) *udppacket.Packet {
	utNack.LOG.Println("ut_nack_module ReadFromServerProxy SendMap size:", len(utNack.SendMap), utNack.minSendId, utNack.maxSendId, "RecvMap size:", len(utNack.RecvMap), utNack.minRecvId, utNack.maxRecvId)
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	utNack.processRecvPacket(p)
	return p
}

func (utNack *UT_NACK) ReadFromClientProxy(p *udppacket.Packet) *udppacket.Packet {
	utNack.LOG.Println("ut_nack_module ReadFromClientProxy SendMap size:", len(utNack.SendMap), utNack.minSendId, utNack.maxSendId, "RecvMap size:", len(utNack.RecvMap), utNack.minRecvId, utNack.maxRecvId)
	utNack.lock.Lock()
    defer utNack.lock.Unlock()
	utNack.processRecvPacket(p)
	return p
}

/************** self function ****************/


/**********************************************************
 * 返回接收队列中的需要确认的包，但是需要考虑一些比较详细的情况，需要wrapper辅助判断
 * - 对于还从来没有确认过的包，直接返回包的tunnelId
 * - 对于已经重传过的包，如果还在重传时效期内，就继续往后推，但是需要返回正确的包类型
 **********************************************************/
func (utNack *UT_NACK) getNackPacketTunnelId() (uint64, byte) {
    utNack.LOG.Println("getNackPacketTunnelId minRecvId:", utNack.minRecvId, "maxRecvId", utNack.maxRecvId, "recv len:", len(utNack.RecvMap))
    flag := udppacket.PACK_ACK
    i := utNack.minRecvId
    for ; i <= utNack.maxRecvId; i++ {
        //utNack.LOG.Println(i)
        if wrapper, ok := utNack.RecvMap[i]; !ok {       // 甚至连wrapper都没有创建，肯定需要nack
            wrapper := new(PACKET_WRAPPER)
            wrapper.status = WRAPPER_STATUS_NACKED
            return i, flag
        } else {
            if wrapper.status == WRAPPER_STATUS_NACKED {
                if wrapper.nackTime.Before(time.Now()) { // 这个包nack超时，重新nack
                    wrapper.nackTime = time.Now().Add(utNack.ackDuration)
                    return i, flag
                } else {                                 // 这个包正在nack，所以去掉ack标志
                    flag = ^udppacket.PACK_ACK & flag 
                    continue
                }
            }
            if wrapper.status == WRAPPER_STATUS_RECVED {  // 正常接收到的包，释放并继续计算
                delete(utNack.RecvMap, i)
                utNack.minRecvId = i + 1
            }
        }
    }
    return i, flag
}
/***********************************************************
 * 处理发送给tunnel的数据包
 * - 功能：将新数据包加入发送队列，重传的数据包，更新状态
 ***********************************************************/
func (utNack *UT_NACK) processSendPacket(p *udppacket.Packet) {
    
    /***** SEND *****/
    if wrapper, ok := utNack.SendMap[p.TunnelId]; ok {
        wrapper.p = p
        wrapper.status = WRAPPER_STATUS_RTR
    } else { // 新包，发送队列中不存在该packet的wrapper
        wrapper = new(PACKET_WRAPPER)
        wrapper.p = p
        wrapper.status = WRAPPER_STATUS_SENDED
        utNack.SendMap[p.TunnelId] = wrapper
    }
    
    if p.TunnelId > utNack.maxSendId {
        utNack.maxSendId = p.TunnelId
    }
    
    /***** RECV *****/
    otherTunnelId, flag := utNack.getNackPacketTunnelId()
    p.ChangeOtherTunnelId(otherTunnelId)
    p.ChangePacketType(flag)
}

/*************************************************************
 * 处理从tunnel收到的数据包
 * - 功能:根据包的类型，处理相应的OtherTunnelId的数据包，是还没有到
 * 该数据包，还是需要重传顺便对需要确认的数据包进行确认释放
 ************************************************************/
func (utNack *UT_NACK) processRecvPacket(p *udppacket.Packet) {
    utNack.LOG.Println("processRecvPacket")
    status := p.GetPacketType()
    utNack.LOG.Println("tunnelId", p.TunnelId, "packetType", status)
    /***** SEND *****/
    if status & udppacket.PACK_ACK == udppacket.PACK_ACK {
        utNack.LOG.Println("ack")
        utNack.ackPackets(p.OtherTunnelId)
    }
     
    /***** RECV *****/
    wrapper, ok := utNack.RecvMap[p.TunnelId]
    if ok {
        wrapper.status = WRAPPER_STATUS_RECVED
    } else {
        wrapper = new(PACKET_WRAPPER)
        wrapper.status = WRAPPER_STATUS_RECVED
    }
    
    utNack.RecvMap[p.TunnelId] = wrapper
    if p.TunnelId > utNack.maxRecvId {
        utNack.maxRecvId = p.TunnelId
    }   
}

func (utNack *UT_NACK) ackPackets(nackTunnelId uint64) {
    for i := utNack.minSendId; i < nackTunnelId; i++ {
        if _, ok := utNack.SendMap[i]; ok {
            delete(utNack.SendMap, i)
            utNack.minSendId = i + 1
        } else {
            utNack.LOG.Print("ackPackets error:missing send packet to ack")
        }
    }
}

/**
 * 修改状态是可行的，因为每次只有一个调用。
 *******/
func (utNack *UT_NACK) nackPacket(tunnelId uint64) *udppacket.Packet {
    return nil
}
