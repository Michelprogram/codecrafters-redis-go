package redis

import (
	"context"
	"encoding/hex"
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

func (_ Ping) IsWritable() bool {
	return false
}

type Echo struct{}

func (_ Echo) Send(conn net.Conn, args [][]byte, _ *Redis) error {

	_, err := conn.Write(createBulkString(args[0]))
	return err
}

func (_ Echo) IsWritable() bool {
	return false
}

type Set struct{}

func (s Set) Send(conn net.Conn, args [][]byte, server *Redis) error {

	var err error

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

	if server.IsMaster {
		_, err = conn.Write(OK)
	}

	return err
}

func (s Set) IsWritable() bool {
	return true
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

func (_ Get) IsWritable() bool {
	return false
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

func (_ Info) IsWritable() bool {
	return false
}

type ReplConf struct {
}

func (_ ReplConf) Send(conn net.Conn, args [][]byte, server *Redis) error {

	var err error

	key := strings.ToLower(string(args[0]))
	
	switch key {
	case "listening-port":
		server.Replications = append(server.Replications, conn)
	}

	_, err = conn.Write(OK)

	return err
}

func (_ ReplConf) IsWritable() bool {
	return false
}

type Psync struct {
}

func (_ Psync) Send(conn net.Conn, _ [][]byte, server *Redis) error {

	data := fmt.Sprintf("+FULLRESYNC %s %d\r\n", server.MasterReplicationId, server.MasterReplicationOffset)

	_, err := conn.Write([]byte(data))

	if err != nil {
		return err
	}

	rdb, err := hex.DecodeString(server.RDB)

	if err != nil {
		return err
	}

	data = fmt.Sprintf("$%d\r\n%s", len(rdb), rdb)

	_, err = conn.Write([]byte(data))

	if err != nil {
		return err
	}

	return err
}

func (_ Psync) IsWritable() bool {
	return false
}
