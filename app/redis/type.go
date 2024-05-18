package redis

import (
	"context"
	"net"
)

type Response []byte

type Data struct {
	Content string
	context.Context
}

type IServer interface {
	ListenAndServe() error
}

type ICommand interface {
	Send(conn net.Conn, args [][]byte, server *Redis) error
	IsWritable() bool
}
