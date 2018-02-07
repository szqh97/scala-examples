package TT

import (
	"../IM/IM_BaseDefine"
	"../IM/IM_Buddy"
	"../IM/IM_Login"
	"github.com/golang/protobuf/proto"
	"log"
)

func (client *ClientConn) handlePdu(pdu *ImPdu) {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Printf("entering handlePdu..., pdu is %v\n", pdu)
	switch pdu.command_id {
	case uint16(IM_BaseDefine.LoginCmdID_CID_LOGIN_RES_USERLOGIN):
		msg := &IM_Login.IMLoginRes{}
		proto.Unmarshal(pdu.msg, msg)
		log.Println(msg)

		client.UserId = msg.GetUserInfo().GetUserId()
		client.IsLogin = true
		log.Println(msg.GetUserInfo().GetCustomerId())
	case uint16(IM_BaseDefine.BuddyListCmdID_CID_BUDDY_LIST_RECENT_CONTACT_SESSION_RESPONSE):
		msg := &IM_Buddy.IMRecentContactSessionRsp{}
		proto.Unmarshal(pdu.msg, msg)
		log.Println(msg.String())

	default:
		log.Fatal("Invalid commdd id ", pdu.command_id)
	}
}
