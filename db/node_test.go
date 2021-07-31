package db

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db OPDB

func InitDB() (OPDB, error) {
	mysqlAddr := "127.0.0.1"
	mysqlPort := "3306"
	mysqlOpDB := "raft"
	mysqlUsr := "raft"
	mysqUsrPwd := "raft"

	return NewDBConfiguration(mysqlUsr, mysqUsrPwd, "mysql", mysqlOpDB, mysqlPort, mysqlAddr).NewDBConnection()
}

func init() {
	var err error
	db, err = InitDB()
	if err != nil {
		panic(err)
	}
}

func TestInsertNodeTable(t *testing.T) {
	port := 1234
	addr := "127.0.0.1:1231"
	id, err := db.InsertDbRecord(port, addr)
	assert.Nil(t, err)
	fmt.Println(id)

	linkUrl, err := db.GetClusterIps()
	assert.Nil(t, err)
	url := strings.Join(linkUrl, ",")

	fmt.Println(url)
}
