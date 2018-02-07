package DBManager

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

const (
	//mycli -utest -p 123456 -h120.26.137.224 -P 23306 teamtalk
	mysql_url string = "test:123456@tcp(120.26.137.224:23306)/teamtalk"
)

var m *DBManager
var once sync.Once

func GetInstance() *DBManager {
	once.Do(func() {
		m = NewDBManager()
	})

	return m
}

type DBManager struct {
	mysqlconn *sql.DB
}

func NewDBManager() *DBManager {
	m := &DBManager{}
	m.Init()
	return m

}

func (p *DBManager) Init() {
	var err error
	p.mysqlconn, err = sql.Open("mysql", mysql_url)
	if err != nil {
		panic(err.Error())
	}
	p.mysqlconn.SetMaxIdleConns(50)
	p.mysqlconn.SetMaxOpenConns(20)
	err = p.mysqlconn.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func (p *DBManager) QueryUserIdByName(Name string) (uint32, error) {
	sql := "select id from IMUser where name = '" + Name + "'"

	log.Printf("sql is [%s]", sql)
	m := GetInstance()
	rows, err := m.mysqlconn.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	var userid uint32
	for rows.Next() {
		err = rows.Scan(&userid)
		fmt.Println(userid)
	}

	return userid, nil
}
