package DBManager

import (
	"testing"
)

func TestQueryUserIdByName(t *testing.T) {
	m := GetInstance()
	uid, err := m.QueryUserIdByName("dj352801")
	if err != nil {
		t.Error(err)
	}
	t.Log(uid)
	if uid != 2396 {
		t.Error("query db failed")
	}

}
