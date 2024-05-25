package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/redis/database"
	"net"
)

type RESP []byte

type Node interface {
	GetDatabase() *database.Database
	AddReplication(conn net.Conn)
	GetInformation() string
	GetRDB() string
	GetMasterInformation() string
	GetOffset() int
	IsMaster() bool
}

type ICommand interface {
	Receive(conn net.Conn, args [][]byte, server Node) error
	IsWritable() bool
}

var (
	BULK_STRING   RESP = []byte("$")
	SIMPLE_STRING RESP = []byte("+")
	ARRAYS        RESP = []byte("*")
	ERROR         RESP = []byte("-")
	CRLF          RESP = []byte("\r\n")
)
