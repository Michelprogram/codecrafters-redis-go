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
	Port         uint
	Address      string
	Commands     map[string]ICommand
	Database     map[string]Data
	Replications []string
	RDB          string
	Information
	net.Listener
}

func newRedis(port uint, role string) *Redis {

	address := fmt.Sprintf("0.0.0.0:%d", port)

	return &Redis{
		Port:         port,
		Address:      address,
		Listener:     nil,
		Information:  newInformation(role),
		Database:     make(map[string]Data),
		Replications: make([]string, 0),
		RDB:          "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2",
		Commands: map[string]ICommand{
			"ping":     Ping{},
			"echo":     Echo{},
			"set":      Set{},
			"get":      Get{},
			"info":     Info{},
			"replconf": ReplConf{},
			"psync":    Psync{},
		},
	}
}

func (r *Redis) send(port string, data []byte) {
	tcpServer, _ := net.ResolveTCPAddr(TCP, "localhost:"+port)

	conn, err := net.DialTCP(TCP, nil, tcpServer)

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(data)

	if err != nil {
		log.Println("Failed replication at " + port)
	}
}

func (r *Redis) propagation(data []byte) {

	for _, port := range r.Replications {
		go r.send(port, data)
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

		cmd, ok := r.Commands[arg]

		if ok {
			err = cmd.Send(conn, args[4:], r)
			if err != nil {
				return err
			}
		} else {
			return errors.New("ICommand " + arg + " doesn't exist")

		}

		if r.IsMaster && cmd.IsWritable() {
			r.propagation(buffer)
		}

	}

	return nil

}
