package ImPdu

import (
	"bytes"
	"encoding/binary"
)

type PduHeader struct {
	length     uint32
	version    uint16
	flag       uint16
	service_id uint16
	command_id uint16
	reversed   uint16
}

type ImPdu struct {
	PduHeader
	msg []byte
}

func (pdu *ImPdu) SetFlag(flag uint16) {
	pdu.flag = flag
}

func (pdu *ImPdu) SetCommandId(cid uint16) {
	pdu.command_id = cid
}

func (pdu *ImPdu) SetServiceId(sid uint16) {
	pdu.service_id = sid
}

func (pdu *ImPdu) SetReversed(reversed uint16) {
	pdu.reversed = reversed
}

func (pdu *ImPdu) SetMsg(msg []byte) {
	pdu.msg = msg
	pdu.length = pdu.length + uint32(len(msg))
}

func (pdu *ImPdu) GetCommandId() uint16 {
	return pdu.command_id
}

func (pdu *ImPdu) GetServiceId() uint16 {
	return pdu.service_id
}

func (pdu *ImPdu) GetMsg() []byte {
	return pdu.msg
}

func (pdu *ImPdu) SerializedToBytes() []byte {
	buf := new(bytes.Buffer)
	codec := binary.LittleEndian
	binary.Write(buf, codec, pdu.length)
	binary.Write(buf, codec, pdu.version)
	binary.Write(buf, codec, pdu.flag)
	binary.Write(buf, codec, pdu.service_id)
	binary.Write(buf, codec, pdu.command_id)
	binary.Write(buf, codec, pdu.reversed)
	binary.Write(buf, codec, pdu.msg)
	return buf.Bytes()
}

func ParserFromBytes(pdu_buf []byte) *ImPdu {
	var pdu ImPdu
	buf := bytes.NewReader(pdu_buf)
	binary.Read(buf, binary.LittleEndian, &pdu)
	return &pdu

}
