package udptunnel

import (
	"udppacket"
	"log"
	"sync"
)

type Pacing struct {
    Index  int
	D int
	lock   *sync.Mutex
    Ut     *UDPTunnel
    LOG    *log.Logger
}

func NewPacing(LOG *log.Logger) *Pacing {
	p := new(Pacing)
    p.LOG = LOG
	return p
}

func (pacing *Pacing)InitHandler(ut *UDPTunnel, index int) {
	pacing.D = 1000
	pacing.lock = new(sync.Mutex)
	pacing.LOG.Println("InitHandler")
    pacing.Ut = ut
    pacing.Index = index
}

func (pacing *Pacing)Debug() string {
	return string(pacing.D)
}

func (pacing *Pacing)WriteToServerProxy(p *udppacket.Packet) *udppacket.Packet {
	pacing.LOG.Println("ut_pacing_module WriteToServerProxy")
	pacing.lock.Lock()
	defer pacing.lock.Unlock()
	return p
}

func (pacing *Pacing)WriteToClientProxy(p *udppacket.Packet) *udppacket.Packet {
	pacing.LOG.Println("ut_pacing_module WriteToClientProxy")
	pacing.lock.Lock()
	defer pacing.lock.Unlock()
	return p
}

func (pacing *Pacing)ReadFromServerProxy(p *udppacket.Packet) *udppacket.Packet {
	pacing.LOG.Println("ut_pacing_module ReadFromServerProxy")
	pacing.lock.Lock()
	defer pacing.lock.Unlock()
	return p
}

func (pacing *Pacing)ReadFromClientProxy(p *udppacket.Packet) *udppacket.Packet {
	pacing.LOG.Println("ut_pacing_module ReadFromClientProxy")
	pacing.lock.Lock()
	defer pacing.lock.Unlock()
	return p
}