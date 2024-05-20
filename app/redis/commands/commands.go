package commands

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Ping struct{}

func (_ Ping) Receive(conn net.Conn, _ [][]byte, server Node) error {

	var err error

	var resp BuilderRESP

	pong := resp.EncodeAsSimpleString("PONG", SIMPLE_STRING)

	if server.IsMaster() {
		_, err = fmt.Fprint(conn, pong.String())
	}

	return err
}

func (_ Ping) IsWritable() bool {
	return false
}

type Echo struct{}

func (_ Echo) Receive(conn net.Conn, args [][]byte, _ Node) error {

	echo := NewBulkString(args[0])

	_, err := fmt.Fprint(conn, echo)

	return err
}

func (_ Echo) IsWritable() bool {
	return false
}

type Set struct{}

func (s Set) Receive(conn net.Conn, args [][]byte, server Node) error {

	var err error
	var builder BuilderRESP

	key, content := string(args[0]), string(args[2])

	fmt.Println(len(args))

	if len(args) > 4 {

		delay, err := strconv.Atoi(string(args[6]))

		if err != nil {
			return err
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(delay))

		server.GetDatabase().Add(key, content, ctx)

	} else {
		server.GetDatabase().Add(key, content, nil)
	}

	if server.IsMaster() {
		_, err = fmt.Fprintf(conn, builder.Ok().String())
	}

	return err
}

func (s Set) IsWritable() bool {
	return true
}

type Get struct{}

func (_ Get) Receive(conn net.Conn, args [][]byte, server Node) error {

	var builder BuilderRESP

	key := string(args[0])

	val, err := server.GetDatabase().Get(key)

	if err != nil {
		_, err = fmt.Fprintf(conn, builder.Null().String())
		return err
	}

	if val.Context == nil {
		_, err = fmt.Fprintf(conn, NewBulkString(val.Content).String())

	} else {
		select {
		case <-val.Done():
			_, err = fmt.Fprintf(conn, builder.Null().String())
		default:
			_, err = fmt.Fprintf(conn, NewBulkString(val.Content).String())
		}
	}

	return err
}

func (_ Get) IsWritable() bool {
	return false
}

type Info struct{}

func (_ Info) Receive(conn net.Conn, args [][]byte, server Node) error {

	var err error

	key := strings.ToLower(string(args[0]))

	switch key {
	case "replication":
		_, err = fmt.Fprintf(conn, NewBulkString(server.GetInformation()).String())
	}

	return err
}

func (_ Info) IsWritable() bool {
	return false
}

type ReplConf struct {
}

func (_ ReplConf) Receive(conn net.Conn, args [][]byte, server Node) error {

	var err error

	var resp BuilderRESP

	key := strings.ToLower(string(args[0]))

	log.Println(key)

	switch key {
	case "listening-port":
		server.AddReplication(conn)
		_, err = fmt.Fprintf(conn, resp.Ok().String())

	case "getack":
		_, err = fmt.Fprintf(conn, resp.EncodeAsArray("REPLCONF", "ACK", "0").String())

	default:
		_, err = fmt.Fprintf(conn, resp.Ok().String())
	}

	return err
}

func (_ ReplConf) IsWritable() bool {
	return false
}

type Psync struct {
}

func (_ Psync) Receive(conn net.Conn, _ [][]byte, server Node) error {

	data := server.GetMasterInformation()

	_, err := conn.Write([]byte(data))

	if err != nil {
		return err
	}

	rdb, err := hex.DecodeString(server.GetRDB())

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

type RDB struct {
}

func (_ RDB) Receive(conn net.Conn, args [][]byte, _ Node) error {

	var resp BuilderRESP

	_, err := fmt.Fprintf(conn, resp.EncodeAsArray("replconf").String())

	return err
}

func (_ RDB) IsWritable() bool {
	return false
}
