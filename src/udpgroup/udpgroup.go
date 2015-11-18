package udpgroup
import (
	"udppacket"
)


type Group struct {
	id				uint32
	Packets			[]*udppacket.Packet

	maxPackets		int
	index			int			// 在被动接收的组中index是没有任何作用的
}

func CreateNewGroup(id uint32, size int) *Group {
	g := new(Group)
	g.maxPackets = size
	g.index = 0
	g.Packets = make([]*udppacket.Packet, g.maxPackets)
	g.id = id
	return g
}

func (g *Group)AddNewPacketData(rawData []byte, dst string) *udppacket.Packet{
	if g.maxPackets == g.index {
		return nil
	}/*
	p := udppacket.CreateNewPacket(len(g.Packets), rawData, dst)
	g.Packets[g.index] = p
	g.index++*/
	return nil
}

func (g *Group)IsFull() bool {
	if g.index >= g.maxPackets {
		return true
	} else {
		return false
	}
}


func (g *Group)AddNewPacketWithId(p *udppacket.Packet) int {
	return 0
}
