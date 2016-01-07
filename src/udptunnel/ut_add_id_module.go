package udptunnel

import (
    "log"
    "sync"
    "udppacket"
)

type UT_ADD_ID struct {
    Ut     *UDPTunnel
    lock   *sync.Mutex
    LOG    *log.Logger
    Index  int
}

type CTX_UT_ADD_ID struct {
    
}

func NewUtAddId (LOG *log.Logger) *UT_ADD_ID {
    utAddId := new(UT_ADD_ID)
    utAddId.LOG = LOG
    utAddId.lock = new(sync.Mutex)
    return utAddId
}

func (utAddId *UT_ADD_ID) Debug() string {
    return ""
}

func (utAddId *UT_ADD_ID) InitHandler(ut *UDPTunnel, index int) {
    utAddId.LOG.Print("InitHandler")
    utAddId.Ut = ut
    utAddId.Index = index
}

func (utAddId *UT_ADD_ID) WriteToServerProxy(p *udppacket.Packet) *udppacket.Packet {
    utAddId.LOG.Print("ut_add_id_module WriteToServerProxy")
    utAddId.lock.Lock()
    defer utAddId.lock.Unlock()
    if p.ModulesCtx[utAddId.Index] == nil {
        p.ModulesCtx[utAddId.Index] = new(CTX_UT_ADD_ID)
        tunnelId := utAddId.Ut.GetNewTunnelId()
        p.ChangeTunnelId(tunnelId)
    }
    return p
}

func (utAddId *UT_ADD_ID) WriteToClientProxy(p *udppacket.Packet) *udppacket.Packet {
    utAddId.LOG.Print("ut_add_id_module WriteToClientProxy")
    utAddId.lock.Lock()
    defer utAddId.lock.Unlock()
    if p.ModulesCtx[utAddId.Index] == nil {
        p.ModulesCtx[utAddId.Index] = new(CTX_UT_ADD_ID)
        tunnelId := utAddId.Ut.GetNewTunnelId()
        p.ChangeTunnelId(tunnelId)
    }
    return p
}

func (utAddId *UT_ADD_ID) ReadFromServerProxy(p *udppacket.Packet) *udppacket.Packet {
    utAddId.LOG.Print("ut_add_id_module ReadFromServerProxy")
    return p
}

func (utAddId *UT_ADD_ID) ReadFromClientProxy(p *udppacket.Packet) *udppacket.Packet {
    utAddId.LOG.Print("ut_add_id_module ReadFromClientProxy")
    return p
}
 
 