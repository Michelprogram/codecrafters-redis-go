package redis

import "net"

type Response []byte

var (
	OK   Response = []byte("+OK\r\n")
	PONG Response = []byte("+PONG\r\n")
	NULL Response = []byte("$-1\r\n")
)

var (
	SPLITTER = []byte("\r\n")
)

type command interface {
	Send(conn net.Conn, args [][]byte, server *Redis) error
}
