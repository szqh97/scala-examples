package TT

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	msg_server_url string = "http://app-test.kaipao.cc/msg_server"
)

type LoginServerRsp struct {
	Code       int32  `json:code`
	BackupIP   string `json:backupIP`
	Discovery  string `json:discovery`
	MsfsBackup string `json:msfsbackup`
	MsfsPrior  string `json:msfsPrior`
	Msg        string `json:msg`
	Port       string `json:port`
	PriorIP    string `json:priorIP`
}
type WaitGrapWrapper struct {
	sync.WaitGroup
}

func (w *WaitGrapWrapper) Wrap(cb func()) {

	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

func GetMsgServerAddress() (string, string, error) {
	resp, err := http.Get(msg_server_url)
	if err != nil {
		log.Printf("request [%s] failed", msg_server_url)
		log.Fatal(err)

		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	var loginresp LoginServerRsp
	err = json.Unmarshal(body, &loginresp)
	return loginresp.PriorIP, loginresp.Port, nil
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
