package TT

import (
	"../DBManager"
	"crypto/md5"
	"fmt"
	"net"
)

type ClientConn struct {
	UserName string
	Passwd   string
	UserId   uint32
	IsLogin  bool

	conn           *net.TCPConn
	pduReceiveChan chan *ImPdu
	pduSenderChan  chan *ImPdu
	exitChan       chan bool
	waitGroup      WaitGrapWrapper
}

func NewClientConn(MsgAddr, UserName string) *ClientConn {

	o1 := md5.Sum([]byte(UserName))
	o2 := md5.Sum([]byte(fmt.Sprintf("%x", o1)))
	passwd := fmt.Sprintf("%x", o2)
	m := DBManager.NewDBManager()
	userId, _ := m.QueryUserIdByName(UserName)
	addr, err := net.ResolveTCPAddr("tcp4", MsgAddr)
	checkErr(err)
	c, err := net.DialTCP("tcp", nil, addr)
	checkErr(err)
	client := ClientConn{
		UserName:       UserName,
		Passwd:         passwd,
		UserId:         userId,
		conn:           c,
		pduReceiveChan: make(chan *ImPdu),
		pduSenderChan:  make(chan *ImPdu),
		exitChan:       make(chan bool),
	}

	client.waitGroup.Wrap(func() { client.readLoop() })
	client.waitGroup.Wrap(func() { client.writeLoop() })
	client.waitGroup.Wrap(func() { client.handleLoop() })

	return &client
}

func (client *ClientConn) writeLoop() {
	for {
		select {
		case pdu := <-client.pduReceiveChan:

			cnt, err := client.conn.Write(pdu.SerializedToBytes())
			if err != nil {
				fmt.Println("write file failed", cnt, err.Error())
			}
		default:
		}
	}
}

func (client *ClientConn) readLoop() {
	for {
		pdu, err := ImPduReader(client)
		if err != nil {
			return
		}
		client.pduSenderChan <- pdu
	}

}
func (client *ClientConn) handleLoop() {
	for {
		select {
		case pdu := <-client.pduSenderChan:
			client.handlePdu(pdu)
		default:
		}
	}
}
