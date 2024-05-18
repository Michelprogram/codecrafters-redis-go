package redis

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Redis struct {
	Port     uint
	Address  string
	Commands map[string]ICommand
	Database map[string]Data
	Information
	net.Listener
}

func newRedis(port uint, role string) *Redis {

	address := fmt.Sprintf("0.0.0.0:%d", port)

	return &Redis{
		Port:        port,
		Address:     address,
		Listener:    nil,
		Information: newInformation(role, address),
		Database:    make(map[string]Data),
		Commands: map[string]ICommand{
			"ping": Ping{},
			"echo": Echo{},
			"set":  Set{},
			"get":  Get{},
			"info": Info{},
		},
	}
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
			return errors.New("ICommand " + arg + " doesn't exist")
		}
	}

	return nil

}
