package utils
import (
	"bytes"
	"encoding/binary"
)
func Int32ToBytes(ui uint32) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, ui)
	return buf.Bytes()
}

func BytesToInt32(data []byte) uint32 {
	buf := bytes.NewBuffer(data)
	var ui uint32
	binary.Read(buf, binary.BigEndian, &ui)
	return ui
}
func Int16ToBytes(i int) []byte {
	buf := make([]byte, 2)
	buf[1] = byte(i & 0xff)
	buf[0] = byte((i >> 8) & 0xff)
	return buf
}

func BytesToInt16(data []byte) int {

	var i int
	i = 0
	i = int(byte(i) | data[0])
	i = i << 8
	i = int(byte(i) | data[1]) + i
	return i
}
