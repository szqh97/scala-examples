package TT

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

type PduHeader struct {
	length     uint32
	version    uint16
	flag       uint16
	service_id uint16
	command_id uint16
	seq_num    uint16
	reversed   uint16
}

var HeadLen int = binary.Size(PduHeader{})

type ImPdu struct {
	PduHeader
	msg []byte
}

func NewImPdu(sid uint16, cid uint16, msg []byte) *ImPdu {
	msglen := len(msg)
	pdu := ImPdu{
		PduHeader: PduHeader{
			length:     uint32(HeadLen + msglen),
			version:    1,
			flag:       0,
			service_id: sid,
			command_id: cid,
			seq_num:    1,
			reversed:   0,
		},
		msg: msg,
	}
	return &pdu
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
	binary.Write(buf, codec, pdu.seq_num)
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

func ImPduReader(client *ClientConn) (*ImPdu, error) {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	fullBuf := bytes.NewBuffer([]byte{})
	var h *PduHeader = nil
	for {
		data := make([]byte, 1024)
		log.Println("before read ....")
		readLen, err := client.conn.Read(data)
		log.Println("after read ....")
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		if readLen == 0 {
			return nil, errors.New("read zeor bytes")
		} else {
			fullBuf.Write(data[:readLen])
			if h == nil && readLen >= 16 {
				hbuf := make([]byte, 16)
				_, _ = fullBuf.Read(hbuf)
				h, err := ParseToHeader(hbuf)
				if err != nil {
					log.Fatal("parse to header failed, ", err.Error())
					return nil, err
				}
				fmt.Println(h)
			} else {
				msgbuf := make([]byte, h.length-16)
				n, err := fullBuf.Read(msgbuf)
				if err != nil {
					return nil, err
				}
				if n != int(h.length-16) {
					return nil, errors.New("msg body truncated!")
				}
				return &ImPdu{
					PduHeader: *h,
					msg:       msgbuf,
				}, nil
			}
		}
	}
}

func ParseToHeader(headerbuf []byte) (*PduHeader, error) {
	if len(headerbuf) != 16 {
		return nil, errors.New("buf length error")
	}
	length_ := uint32(binary.LittleEndian.Uint32(headerbuf[:4]))
	version_ := uint16(binary.LittleEndian.Uint16(headerbuf[4:6]))
	flag_ := uint16(binary.LittleEndian.Uint16(headerbuf[6:8]))
	service_id_ := uint16(binary.LittleEndian.Uint16(headerbuf[8:10]))
	command_id_ := uint16(binary.LittleEndian.Uint16(headerbuf[10:12]))
	seq_num_ := uint16(binary.LittleEndian.Uint16(headerbuf[12:14]))
	reversed_ := uint16(binary.LittleEndian.Uint16(headerbuf[14:16]))
	return &PduHeader{
		length:     length_,
		version:    version_,
		flag:       flag_,
		service_id: service_id_,
		command_id: command_id_,
		seq_num:    seq_num_,
		reversed:   reversed_,
	}, nil
}

func ReadPduFromReader(client *ClientConn) (*ImPdu, error) {
	log.Println("in ReadPduFromReader ...")
	headerSize := binary.Size(PduHeader{})
	headerbuf := make([]byte, 1024)
	log.Printf("kKKKKKKKKKKKK con:%v \n", client.conn)
	cnt, err := client.conn.Read(headerbuf)
	log.Printf("kKKKKKKKKKKKK con:%v \n", client.conn)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	if cnt != headerSize {
		log.Fatal("read %d byte", cnt)
		errmsg := fmt.Sprintf("expect %d bytes but read %d only", headerSize, cnt)
		err := errors.New(errmsg)
		return nil, err
	}

	length_ := uint32(binary.LittleEndian.Uint32(headerbuf[:4]))
	version_ := uint16(binary.LittleEndian.Uint16(headerbuf[4:6]))
	flag_ := uint16(binary.LittleEndian.Uint16(headerbuf[6:8]))
	service_id_ := uint16(binary.LittleEndian.Uint16(headerbuf[8:10]))
	command_id_ := uint16(binary.LittleEndian.Uint16(headerbuf[10:12]))
	seq_num_ := uint16(binary.LittleEndian.Uint16(headerbuf[12:14]))
	reversed_ := uint16(binary.LittleEndian.Uint16(headerbuf[14:16]))
	msg_ := make([]byte, int(length_)-headerSize)
	cnt, err = client.conn.Read(msg_)

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	pdu := ImPdu{
		PduHeader: PduHeader{
			length:     length_,
			version:    version_,
			flag:       flag_,
			service_id: service_id_,
			command_id: command_id_,
			seq_num:    seq_num_,
			reversed:   reversed_,
		},
		msg: msg_,
	}
	return &pdu, nil

}
