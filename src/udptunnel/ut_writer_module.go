package udptunnel

import (
    "log"
    "udppacket"
)

type UT_WRITER struct {
    Index       int
    Ut          *UDPTunnel
    LOG         *log.Logger
}

func NewUtWriter(LOG *log.Logger) *UT_WRITER {
    utWriter := new(UT_WRITER)
    utWriter.LOG = LOG
    return utWriter
}

func (utWruter *UT_WRITER) Debug() string {
    return ""
}

func (utWriter *UT_WRITER) InitHandler(ut *UDPTunnel, index int) {
    utWriter.LOG.Print("InitHandler")
    utWriter.Ut = ut
    utWriter.Index = index
}

func (utWriter *UT_WRITER) WriteToServerProxy(p *udppacket.Packet) *udppacket.Packet {
    utWriter.Ut.WriteData(p.GetPacket())
    return nil
}

func (utWriter *UT_WRITER) WriteToClientProxy(p *udppacket.Packet) *udppacket.Packet {
    utWriter.Ut.WriteData(p.GetPacket())    
    return nil
}

func (utWriter *UT_WRITER) ReadFromServerProxy(p *udppacket.Packet) *udppacket.Packet {
    return p
}

func (utWriter *UT_WRITER) ReadFromClientProxy(p *udppacket.Packet) *udppacket.Packet {
    return p
}