package main

import (
	"./TT"
	"fmt"
	"time"
)

func main() {
	fmt.Println("vim-go")

	ip, port, _ := TT.GetMsgServerAddress()
	msgaddr := ip + ":" + port
	client := TT.NewClientConn(msgaddr, "dj352801")
	client.Login()
	<-time.After(time.Second * 2)
	client.GetRecentSession(0)

	<-time.After(time.Second * 50)
}
