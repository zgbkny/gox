package udptunnel

import (
	"udppacket"
	"log"
	"sync"
)

type Pacing struct {
	D int
	lock *sync.Mutex
}

func NewPacing() *Pacing {
	p := new(Pacing)
	return p
}

func (pacing *Pacing)InitHandler() {
	pacing.D = 1000
	pacing.lock = new(sync.Mutex)
	log.Println("InitHandler", pacing.D)
}

func (pacing *Pacing)WriteToServerProxy(p *udppacket.Packet) bool {
	log.Println("ut_pacing_module WriteToServerProxy")
	pacing.lock.Lock()
	defer pacing.lock.Unlock()
	
	return true
}

func (pacing *Pacing)WriteToClientProxy(p *udppacket.Packet) bool {
	log.Println("ut_pacing_module WriteToClientProxy")
	pacing.lock.Lock()
	defer pacing.lock.Unlock()
	
	return true
}

func (pacing *Pacing)ReadFromServerProxy(p *udppacket.Packet) bool {
	log.Println("ut_pacing_module ReadFromServerProxy")
	pacing.lock.Lock()
	defer pacing.lock.Unlock()
	
	return true
}

func (pacing *Pacing)ReadFromClientProxy(p *udppacket.Packet) bool {
	log.Println("ut_pacing_module ReadFromClientProxy")
	pacing.lock.Lock()
	defer pacing.lock.Unlock()
	
	return true
}