package commands

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/redis/database"
	"log"
	"net"
	"os"
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
	} else {
		log.Println("Secondary node doesn't send back pong")
	}

	return err
}

func (_ Ping) IsWritable() bool {
	return true
}

type Echo struct{}

func (_ Echo) Receive(conn net.Conn, args [][]byte, _ Node) error {

	_, err := fmt.Fprint(conn, NewBulkString(args[0]).String())

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
	} else {
		log.Println("Doesnt send set back as secondary")
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

	if val, _ := server.GetConfiguration("dir"); val != "" {

		dbfilename, _ := server.GetConfiguration("dbfilename")

		path := fmt.Sprintf("%s/%s", val, dbfilename)

		file, err := os.ReadFile(path)

		start := bytes.IndexByte(file, 251)
		end := bytes.IndexByte(file[start:], 255) + start

		line := file[start+1 : end]

		size := int(line[3]) + 4

		_, _ = fmt.Fprintf(os.Stdout, string(line[size:]))

		b := NewBulkString(string(line[size+1:]))

		_, err = fmt.Fprintf(conn, b.String())

		return err

	}

	val, err := server.GetDatabase().Get(key)

	if err != nil {
		_, err = fmt.Fprintf(conn, builder.Null().String())
		return err
	}

	if val.Context == nil {
		_, err = fmt.Fprintf(conn, NewBulkString(val.String()).String())

	} else {
		select {
		case <-val.Done():
			_, err = fmt.Fprintf(conn, builder.Null().String())
		default:
			_, err = fmt.Fprintf(conn, NewBulkString(val.String()).String())
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

		offset := strconv.Itoa(server.GetOffset())

		log.Println("Offset " + offset)

		_, err = fmt.Fprintf(conn, resp.EncodeAsArray("REPLCONF", "ACK", offset).String())
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

	_, err := fmt.Fprintf(conn, resp.EncodeAsArray("replconf", "ACK", "0").String())

	return err
}

func (_ RDB) IsWritable() bool {
	return false
}

type Type struct {
}

func (_ Type) Receive(conn net.Conn, args [][]byte, server Node) error {

	var resp BuilderRESP

	key := string(args[0])

	val, err := server.GetDatabase().Get(key)

	if err != nil {
		_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(val.Type, SIMPLE_STRING).String())
		return err
	}

	if val.Context == nil {
		_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(val.Type, SIMPLE_STRING).String())

	} else {
		select {
		case <-val.Done():
			_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(val.Type, SIMPLE_STRING).String())
		default:
			_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(val.Type, SIMPLE_STRING).String())
		}
	}

	return err
}

func (_ Type) IsWritable() bool {
	return false
}

type Xadd struct {
}

func (_ Xadd) Receive(conn net.Conn, args [][]byte, server Node) error {

	var err error
	var res *database.ID
	var resp BuilderRESP

	key := args[0]

	id := args[2]

	for i := 4; i < len(args)-2; i += 2 {
		res, err = server.GetDatabase().AddX(string(key), string(id), args[i], args[i+2])
		if err != nil {
			log.Println(err)
			_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
			return err
		}
	}

	sequence := fmt.Sprintf("%d-%d", res.MillisecondsTime, res.SequenceNumber)

	_, err = fmt.Fprintf(conn, NewBulkString(sequence).String())

	return err
}

func (_ Xadd) IsWritable() bool {
	return false
}

type Xrange struct {
}

func (_ Xrange) Receive(conn net.Conn, args [][]byte, server Node) error {

	var err error
	var resp BuilderRESP

	key := args[0]

	stream, err := server.GetDatabase().Range(string(key), args[2], args[4])

	if err != nil {
		_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())

		return err
	}

	res := resp.XRange(*stream).String()

	_, err = fmt.Fprintf(conn, res)

	return err
}

func (_ Xrange) IsWritable() bool {
	return false
}

type XRead struct {
}

func (_ XRead) Receive(conn net.Conn, args [][]byte, server Node) error {

	var err error
	var stream *database.Stream
	var resp BuilderRESP
	var ctx context.Context

	//TODO Refacto

	if string(args[0]) == "block" {

		timeout, err := strconv.Atoi(string(args[2]))

		key := args[6]

		if err != nil {
			_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
			return err
		}

		data, err := server.GetDatabase().Get(string(key))

		if err != nil {
			_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
			return err
		}

		subscriber := make(chan *database.Stream)

		stream = data.Content.(*database.Stream)

		stream.AddSubscribe(subscriber)

		if timeout == 0 {
			ctx = context.Background()
		} else {
			ctx, _ = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
		}

		select {
		case <-ctx.Done():
			resp = *resp.Null()
		case stream = <-subscriber:
			resp = *resp.XRead(key, *stream)
		}

	} else if len(args) > 6 {

		streams := make([]*database.Stream, 0)
		keys := make([][]byte, 0)

		for i := 2; i < len(args)/2; i += 2 {
			stream, err = server.GetDatabase().Read(string(args[i]), args[i+4])

			if err != nil {
				_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
				return err
			}

			keys = append(keys, args[i])

			streams = append(streams, stream)
		}

		resp.XReadMultiple(keys, streams)

	} else {
		key := args[2]

		id := args[4]

		stream, err = server.GetDatabase().Read(string(key), id)

		if err != nil {
			_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
			return err
		}

		resp.XRead(key, *stream)

	}

	log.Println(conn)

	log.Println(resp.String())

	_, err = fmt.Fprintf(conn, resp.String())

	return err
}

func (_ XRead) IsWritable() bool {
	return false
}

type Config struct {
}

func (_ Config) Receive(conn net.Conn, args [][]byte, server Node) error {

	var resp BuilderRESP

	key := args[2]

	val, err := server.GetConfiguration(string(key))

	if err != nil {
		_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
		return err
	}

	resp.GetConfig(key, []byte(val))

	_, err = fmt.Fprintf(conn, resp.String())

	return err
}

func (_ Config) IsWritable() bool {
	return false
}

type Keys struct {
}

func (_ Keys) Receive(conn net.Conn, args [][]byte, server Node) error {

	var resp BuilderRESP

	dir, err := server.GetConfiguration("dir")

	if err != nil {
		_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
		return err
	}

	dbfilename, err := server.GetConfiguration("dbfilename")

	if err != nil {
		_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
		return err
	}

	path := fmt.Sprintf("%s/%s", dir, dbfilename)

	file, err := os.ReadFile(path)

	if err != nil {
		_, err = fmt.Fprintf(conn, resp.EncodeAsSimpleString(err.Error(), ERROR).String())
		return err
	}

	start := bytes.IndexByte(file, 251)
	end := bytes.IndexByte(file[start:], 255) + start

	key := file[start+1 : end]

	size := int(key[3]) + 4

	b := resp.EncodeAsArray(string(key[4:size]))

	_, err = fmt.Fprintf(conn, b.String())

	return err
}

func (_ Keys) IsWritable() bool {
	return false
}
