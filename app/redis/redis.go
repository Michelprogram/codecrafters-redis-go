package redis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	TCP = "tcp"
)

type Data struct {
	Content string
	context.Context
}

type Redis struct {
	Port     uint
	Role     string
	Address  string
	Commands map[string]command
	Database map[string]Data
	net.Listener
}

func NewServer(port uint, role string) *Redis {

	if role != "master" {
		role = "slave"
	}

	return &Redis{
		Port:     port,
		Address:  fmt.Sprintf("0.0.0.0:%d", port),
		Listener: nil,
		Role:     role,
		Database: make(map[string]Data),
		Commands: map[string]command{
			"ping": Ping{},
			"echo": Echo{},
			"set":  Set{},
			"get":  Get{},
			"info": Info{},
		},
	}
}

func (r *Redis) ListenAndServe() error {
	l, err := net.Listen(TCP, r.Address)

	r.Listener = l

	if err != nil {
		return err
	}

	defer l.Close()

	r.handleRequests()

	return nil
}

func (r *Redis) handleRequests() {
	for {

		conn, err := r.Accept()
		if err != nil {
			log.Printf("Can't handle request : %s\n", err)
		}

		go func() {
			err := r.response(conn)
			if err != nil {
				log.Printf("Can't send response : %s\n", err)
			}
		}()

	}
}

func (r *Redis) response(conn net.Conn) error {

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)

		size, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		args := bytes.Split(buffer[:size], SPLITTER)

		arg := string(bytes.ToLower(args[2]))

		log.Printf("Command received : %s\n", arg)

		if val, ok := r.Commands[arg]; ok {
			err = val.Send(conn, args[4:], r)
			if err != nil {
				return err
			}

		} else {
			return errors.New("command " + arg + " doesn't exist")
		}
	}

	return nil

}
