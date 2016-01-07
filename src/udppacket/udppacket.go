package udppacket
import ( 
	"log"
	"utils"
)


const Length_INDEX = 0
const Length_END = 2
const TunnelId_INDEX = 2
const TunnelId_END = 10
const SessionId_INDEX = 10
const SessionId_END = 14
const Id_INDEX = 14
const Id_END = 18
const ProtoType_INDEX = 18
const PacketType_INDEX = 19
const OtherTunnelId_INDEX = 20
const OtherTunnelId_END = 28
const OtherLen_INDEX = 29


const PACK_NEW byte = 0x01      // 刚刚创建的新包，各个模块需要对该包进行完善
const PACK_NORM byte = 0x02     // 已经不是一个新包了
const PACK_NACK byte = 0x04     // nack
const PACK_ACK byte = 0x08      // ack

/***********************************************************************************
 *
 *
 *  |--Length(2)----|----------TunnelId(8)---------|--------SessionId(4)---------|
 *  |-------------Id(4)-----------|Proto(1)|Pack(1)|---------TunnelId(8)---------|
 *  |Other(1)|--option(n)--|----------------------data---------------------------|
 *
 *  Pack：数据包的标志位，一共有8个标志位
 *  |7|6|5|4|3|2|1|-0(NACK)-|
 *
 *
 ***********************************************************************************/

type Packet struct {
    LOG             *log.Logger
    ModulesCtx      []interface{}
	RawData			[]byte
	Start			int 		// 主动创建的包，Start代表头部开始的位置；从数据解析出的包，Start代表包体开始的位置
    Status          int         // 包状态，比如：重传？还是新包

	next			*Packet
	
	// 控制信息
	Acked			bool
	/******包信息*******/
	Length			int
	TunnelId		uint64
	SessionId		uint32
	Id				uint32
	ProtoType		byte
	PacketType		byte
    OtherTunnelId   uint64
	OtherLen		int
	Dst				string
}

func CreateNewPacket(id uint32, rawData []byte, dst string, modulesCount int, LOG *log.Logger) *Packet {
	//LOG.Println("udppacket createNewPacket")
	p := new(Packet)
    p.ModulesCtx = make([]interface{}, modulesCount)
	p.RawData = rawData
	p.SessionId = 0
	p.TunnelId = 0
	p.Id = id
	p.Dst = dst
	p.Length = len(rawData) - 96
    p.LOG = LOG
    p.Start = 0
	return p
}

func GenPacketFromData(data []byte, LOG *log.Logger) *Packet {

	dataLen := utils.BytesToInt16(data[Length_INDEX : Length_END])
	tunnelId := utils.BytesToUint64(data[TunnelId_INDEX : TunnelId_END])
	sessionId := utils.BytesToUint32(data[SessionId_INDEX : SessionId_END])
	id := utils.BytesToUint32(data[Id_INDEX : Id_END])
	protocal := data[ProtoType_INDEX]
	packetType := data[PacketType_INDEX]
    otherTunnelId := utils.BytesToUint64(data[OtherTunnelId_INDEX : OtherTunnelId_END])
	otherLen := int(data[OtherLen_INDEX])
	start := OtherLen_INDEX + 1
	dst := ""
	if otherLen != 0 {
		log.Println("dst len", otherLen)
		start += int(otherLen)
		dst = string(data[OtherLen_INDEX + 1:otherLen + OtherLen_INDEX + 1])
	}
	p := new(Packet)
	p.Length = dataLen
	p.TunnelId = tunnelId
	p.SessionId = sessionId
	p.Id = id
	p.ProtoType = protocal
	p.PacketType = packetType
    p.OtherTunnelId = otherTunnelId
	p.Start = start
	p.Dst = dst
	p.RawData = data
    p.LOG = LOG
	return p
}

func (p *Packet)LogPacket() {
	//p.LOG.Println("Length", p.Length, "TunnelId", p.TunnelId,  "SessionId", p.SessionId, "Id", p.Id, "ProtoType", p.ProtoType, "PacketType", p.PacketType, "otherlen", p.OtherLen, "Dst", p.Dst)
}

func (p *Packet)genHeader() []byte {
	header := make([]byte, 96)
	dataLenBytes := utils.Int16ToBytes(len(p.RawData) - 96)
	copy(header[Length_INDEX:Length_END], dataLenBytes)

	tunnelIdBytes := utils.Uint64ToBytes(p.TunnelId)
	copy(header[TunnelId_INDEX:TunnelId_END], tunnelIdBytes)

	sessionIdBytes := utils.Uint32ToBytes(p.SessionId)
	copy(header[SessionId_INDEX:SessionId_END], sessionIdBytes)

	packetIdBytes := utils.Uint32ToBytes(p.Id)
	copy(header[Id_INDEX:Id_END], packetIdBytes)

	// 传输层协议类型
	header[ProtoType_INDEX] = p.ProtoType

	// 包类型
	header[PacketType_INDEX] = p.PacketType

	// other 长度
	header[OtherLen_INDEX] = byte(p.OtherLen)
	count := OtherLen_INDEX + 1
	if p.OtherLen != 0 && p.OtherLen == len(p.Dst){
		dstBytes := []byte(p.Dst)
		copy(header[OtherLen_INDEX + 1:OtherLen_INDEX + 1 + p.OtherLen], dstBytes)
		count += len(dstBytes)
	} else {
		//p.LOG.Println("error udppacket genHeader OtherLen not match len(Dst)")
		header[OtherLen_INDEX] = byte(0)
	}
	return header[:count]
}

func (p *Packet)RawDataAddHeader() {
	//p.LOG.Println("udpapcket rawDataAddHeader")
	header := p.genHeader() 
	p.Start = 96 - len(header)
	copy(p.RawData[p.Start:96], header)
}

func (p *Packet)GetPacket() []byte {
	return p.RawData[p.Start:]
}

func (p *Packet)ChangeTunnelId(tunnelId uint64) {
	p.TunnelId = tunnelId
	data := utils.Uint64ToBytes(tunnelId)
	for i := TunnelId_INDEX; i < TunnelId_END; i++ {
		p.RawData[p.Start + i] = data[i - TunnelId_INDEX]
	} 
}

func (p *Packet)ChangeOtherTunnelId(tunnelId uint64) {
    p.OtherTunnelId = tunnelId
    data := utils.Uint64ToBytes(tunnelId)
    for i := OtherTunnelId_INDEX; i < OtherTunnelId_END; i++ {
        p.RawData[p.Start + i] = data[i - OtherTunnelId_INDEX]
    }
}

func (p *Packet)GetPacketType() byte {
    //p.LOG.Println("RawData len:", len(p.RawData), "index:", p.Start + PacketType_INDEX)
    return p.PacketType
}

func (p *Packet)ChangePacketType(packetType byte){
    p.RawData[p.Start + PacketType_INDEX] = p.RawData[p.Start + PacketType_INDEX] | packetType
}

func (p *Packet)SetPacket(packetFlag byte) {
    
}

func (p *Packet)ChangeGroupId(groupId uint32) {

}
