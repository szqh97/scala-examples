package TT

import (
	"../DBManager"
	"../IM/IM_BaseDefine"
	"../IM/IM_Login"
	"crypto/md5"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"time"
)

type ClientConn struct {
	UserName string
	Passwd   string
	UserId   uint32

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
	c.SetReadDeadline(time.Now().Add(time.Second * 20))
	c.SetWriteDeadline(time.Now().Add(time.Second * 20))
	client := ClientConn{
		UserName:       UserName,
		Passwd:         passwd,
		UserId:         userId,
		conn:           c,
		pduReceiveChan: make(chan *ImPdu),
		pduSenderChan:  make(chan *ImPdu),
		exitChan:       make(chan bool),
	}

	client.conn.SetKeepAlive(true)
	client.conn.SetNoDelay(true)
	//client.waitGroup.Wrap(func() { client.IOLoop() })
	client.waitGroup.Wrap(func() { client.readLoop() })
	client.waitGroup.Wrap(func() { client.writeLoop() })
	client.waitGroup.Wrap(func() { client.handleLoop() })

	return &client

}

func (client *ClientConn) writeLoop() {
	for {
		select {
		case pdu := <-client.pduReceiveChan:
			client.conn.Write(pdu.SerializedToBytes())
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
		}
	}
}

func (client *ClientConn) IOLoop() {
	log.Println("entering client IOLoop...")
	go func() {
		for {
			select {
			case requestPdu := <-client.pduReceiveChan:
				_, err := client.conn.Write(requestPdu.SerializedToBytes())
				log.Println(" received from chan", requestPdu.SerializedToBytes())
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}

	}()
	for {

		//	pdu, err := ReadPduFromReader(client)
		pdu, err := ImPduReader(client)
		if err != nil {
			log.Panic(err.Error())
		}

		fmt.Println("read from conn", pdu)

		client.handlePdu(pdu)
		log.Println("xxxxxxx")
	}
	select {
	case <-client.exitChan:
		log.Println("exiting loop")
	}
	log.Println("leaving IOLOOP >>>>>>")

}

func (client *ClientConn) request(pdu *ImPdu) {
	log.Println(" request ...")
	client.pduReceiveChan <- pdu
}

func (client *ClientConn) handlePdu(pdu *ImPdu) {
	log.Println("entering handlePdu...")
	switch pdu.command_id {
	case uint16(IM_BaseDefine.LoginCmdID_CID_LOGIN_RES_USERLOGIN):
		log.Printf("login response :%v\n", *pdu)
	}
}

func (client *ClientConn) Login() {
	msg := IM_Login.IMLoginReq{}
	msg.UserName = &client.UserName
	msg.Password = &client.Passwd
	online := IM_BaseDefine.UserStatType_USER_STATUS_ONLINE
	w := IM_BaseDefine.ClientType_CLIENT_TYPE_WINDOWS
	msg.OnlineStatus = &online
	msg.ClientType = &w

	out, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
	}
	pdu := NewImPdu(uint16(IM_BaseDefine.ServiceID_SID_LOGIN),
		uint16(IM_BaseDefine.LoginCmdID_CID_LOGIN_REQ_USERLOGIN), out)
	log.Println("login request...")
	client.request(pdu)
}
