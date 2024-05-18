package redis

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func createBulkString[K []byte | string](data K) []byte {
	res := fmt.Sprintf("$%d\r\n%s\r\n", len(data), data)

	return []byte(res)
}

type Ping struct{}

func (_ Ping) Send(conn net.Conn, _ [][]byte, _ *Redis) error {
	_, err := conn.Write(PONG)

	return err
}

type Echo struct{}

func (_ Echo) Send(conn net.Conn, args [][]byte, _ *Redis) error {

	_, err := conn.Write(createBulkString(args[0]))
	return err
}

type Set struct{}

func (_ Set) Send(conn net.Conn, args [][]byte, server *Redis) error {

	key, content := string(args[0]), string(args[2])

	if len(args) > 4 {

		delay, err := strconv.Atoi(string(args[6]))

		if err != nil {
			return err
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(delay))

		server.Database[key] = Data{
			Content: content,
			Context: ctx,
		}
	} else {
		server.Database[key] = Data{
			Content: content,
		}
	}

	_, err := conn.Write(OK)

	return err
}

type Get struct{}

func (_ Get) Send(conn net.Conn, args [][]byte, server *Redis) error {

	key := string(args[0])
	val, ok := server.Database[key]

	var err error

	if ok {
		if val.Context == nil {
			_, err = conn.Write(createBulkString(val.Content))
		} else {
			select {
			case <-val.Done():
				_, err = conn.Write(NULL)
			default:
				_, err = conn.Write(createBulkString(val.Content))
			}
		}
	} else {
		_, err = conn.Write(NULL)
	}

	return err
}

type Info struct{}

func (_ Info) Send(conn net.Conn, args [][]byte, server *Redis) error {

	var err error

	key := strings.ToLower(string(args[0]))

	switch key {
	case "replication":
		_, err = conn.Write(createBulkString(server.Information.String()))
	}

	return err
}
