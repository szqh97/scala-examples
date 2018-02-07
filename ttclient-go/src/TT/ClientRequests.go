package TT

import (
	"../IM/IM_BaseDefine"
	"../IM/IM_Buddy"
	"../IM/IM_Login"
	"github.com/golang/protobuf/proto"
	"log"
)

func (client *ClientConn) request(pdu *ImPdu) {
	client.pduReceiveChan <- pdu
}

func (client *ClientConn) Login() {
	msg := IM_Login.IMLoginReq{}
	msg.UserName = &client.UserName
	msg.Password = &client.Passwd
	online := IM_BaseDefine.UserStatType_USER_STATUS_ONLINE
	w := IM_BaseDefine.ClientType_CLIENT_TYPE_ANDROID
	version := "v1.1.0"
	msg.OnlineStatus = &online
	msg.ClientType = &w
	msg.ClientVersion = &version
	out, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
	}
	pdu := NewImPdu(uint16(IM_BaseDefine.ServiceID_SID_LOGIN),
		uint16(IM_BaseDefine.LoginCmdID_CID_LOGIN_REQ_USERLOGIN), out)
	log.Println("login request...")

	client.pduReceiveChan <- pdu
}

func (client *ClientConn) GetRecentSession(ts uint32) {
	msg := IM_Buddy.IMRecentContactSessionReq{
		UserId:           &client.UserId,
		LatestUpdateTime: &ts,
	}
	out, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatal(err)
		return
	}
	pdu := NewImPdu(uint16(IM_BaseDefine.ServiceID_SID_BUDDY_LIST),
		uint16(IM_BaseDefine.BuddyListCmdID_CID_BUDDY_LIST_RECENT_CONTACT_SESSION_REQUEST), out)
	log.Println("get recent session")
	client.pduReceiveChan <- pdu
}
