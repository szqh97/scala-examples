package TT

import "testing"

func TestGetMsgServerAddress(t *testing.T) {
	ip, port, err := GetMsgServerAddress()
	if err != nil {
		t.Log(err)
	}
	t.Log(ip, port)

}
