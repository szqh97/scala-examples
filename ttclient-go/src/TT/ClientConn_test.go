package TT

import (
	"../IM/IM_BaseDefine"
	"testing"
)

func Test_NewClientConn(t *testing.T) {
	client := NewClientConn("120.26.137.224:28000", "dj352801")
	t.Log("%v\n", client)
	t.Log("cmd id %d\n", IM_BaseDefine.LoginCmdID_CID_LOGIN_REQ_DEVICETOKEN)
	client.Login()
}

func Test_Login(t *testing.T) {

}
