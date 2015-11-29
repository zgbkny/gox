package udppacket
import ( 
	"log"
	"utils"
)


const PACK_NORM byte = 0
const PACK_ACK = 1
const PACK_RETR = 2
const PACK_CLOSE = 3

type Packet struct {
	RawData			[]byte
	Start			int

	// 控制信息
	Acked			bool
	/******包信息*******/
	Length			int
	SessionId		uint32
	GroupId			uint32
	Id				uint32
	ProtoType		byte
	PacketType		byte
	Dst				string

}

func CreateNewPacket(id uint32, rawData []byte, dst string) *Packet {
	log.Println("udppacket createNewPacket")
	p := new(Packet)
	p.RawData = rawData
	p.SessionId = 0
	p.GroupId = 0
	p.Id = id 
	p.Dst = dst
	p.Length = len(rawData) - 96
	return p
}

func GenPacketFromData(data []byte) *Packet {

	dataLen := utils.BytesToInt16(data[0:2])
	sessionId := utils.BytesToInt32(data[2:6])
	id := utils.BytesToInt32(data[6:10])
	protocal := data[10]
	packetType := data[11]
	otherLen := int(data[12])
	start := 13
	dst := ""
	if otherLen != 0 {
		log.Println("dst len", otherLen)
		start += int(otherLen)
		dst = string(data[13:otherLen + 13])
	}
	p := new(Packet)
	p.Length = dataLen
	p.SessionId = sessionId
	p.Id = id
	p.ProtoType = protocal
	p.PacketType = packetType
	p.Start = start
	p.Dst = dst
	p.RawData = data
	return p
}

func (p *Packet)LogPacket() {
	log.Println("Length", p.Length, "SessionId", p.SessionId, "GroupId", p.GroupId, "Id", p.Id, "ProtoType", p.ProtoType, "PacketType", p.PacketType, "Dst", p.Dst)
}

func (p *Packet)RawDataAddHeader(header []byte) {
	log.Println("udpapcket rawDataAddHeader")
	p.Start = 96 - len(header)
	copy(p.RawData[p.Start:96], header)
}
func (p *Packet)GetPacket() []byte {
	return p.RawData[p.Start:]
}
