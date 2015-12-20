package udptunnel

import (
	"log"
	"sync"
	"udppacket"
)

type UT_NACK struct {
	D 		int
	lock 	*sync.Mutex
	
}

func NewNack() *UT_NACK {
	utNack := new(UT_NACK)
	return utNack
}

func (utNack *UT_NACK)Debug() string {
	return string(utNack.D)
}

func (utNack *UT_NACK)InitHandler() {
	utNack.D = 100345
	utNack.lock = new(sync.Mutex)
	log.Println("InitHandler", utNack.D)
}

func (utNack *UT_NACK)WriteToServerProxy(p *udppacket.Packet) *udppacket.Packet {
	log.Println("ut_nack_module WriteToServerProxy")
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	
	return p
}

func (utNack *UT_NACK)WriteToClientProxy(p *udppacket.Packet) *udppacket.Packet {
	log.Println("ut_nack_module WriteToClientProxy")
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	
	return p
}

func (utNack *UT_NACK)ReadFromServerProxy(p *udppacket.Packet) *udppacket.Packet {
	log.Println("ut_nack_module ReadFromServerProxy")
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	
	return p
}

func (utNack *UT_NACK)ReadFromClientProxy(p *udppacket.Packet) *udppacket.Packet {
	log.Println("ut_nack_module ReadFromClientProxy")
	utNack.lock.Lock()
	defer utNack.lock.Unlock()
	
	return p
}

