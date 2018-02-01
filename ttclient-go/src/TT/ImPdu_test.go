package TT

import (
	"../IM/IM_BaseDefine"
	"testing"
)

func Test_dumpImPdu(t *testing.T) {
	p := new(ImPdu)
	p.SetFlag(1)

	t.Log("fsdf")
	t.Log(IM_BaseDefine.LoginCmdID_CID_LOGIN_REQ_DEVICETOKEN)
	t.Logf("%v", p.SerializedToBytes())
}
