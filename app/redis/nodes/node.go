package nodes

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/redis/commands"
	"github.com/codecrafters-io/redis-starter-go/app/redis/database"
	"io"
	"log"
	"net"
)

type Node struct {
	Port         uint
	Address      string
	RDB          string
	IsPrimary    bool
	Replications []net.Conn
	Offset       int
	Information
	net.Listener
	database.Database
	commands.Parser
}

func NewNode(port uint, role, dir, dbfilename string) IServer {

	if role == "" {
		return NewPrimary(port, role, dir, dbfilename)
	}

	return NewSecondary(port, role, dir, dbfilename)

}

func (r *Node) propagation(data []byte) {

	for _, replication := range r.Replications {

		data = bytes.Replace(data, []byte("\x00"), []byte(""), -1)

		_, err := replication.Write(data)

		if err != nil {
			log.Printf("Couldn't write to %s : %s \n", replication, err)
		}
	}

}

func (r *Node) handleRequests() {

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

func (r *Node) response(conn net.Conn) error {

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

		args := bytes.Split(buffer[:size], []byte("\r\n"))

		arg := string(bytes.ToLower(args[2]))

		cmd, ok := r.Commands[arg]

		if ok {
			err = cmd.Receive(conn, args[4:], r)
			r.Offset += size
			if err != nil {
				log.Println(err)
			}
		} else {
			return errors.New("ICommand " + arg + " doesn't exist")
		}

		log.Println("Command sended" + arg)

		if r.IsPrimary {
			r.propagation(buffer)
		}

	}

	return nil

}

func (r *Node) AddReplication(conn net.Conn) {
	r.Replications = append(r.Replications, conn)
}

func (r *Node) GetDatabase() *database.Database {
	return &r.Database
}

func (r *Node) IsMaster() bool {
	return r.IsPrimary
}

func (r *Node) GetInformation() string {
	return r.Information.String()
}

func (r *Node) GetRDB() string {
	return r.RDB
}

func (r *Node) GetMasterInformation() string {
	return fmt.Sprintf("+FULLRESYNC %s %d\r\n", r.MasterReplicationId, r.MasterReplicationOffset)
}

func (r *Node) GetOffset() int {
	return r.Offset
}

func (r *Node) GetConfiguration(key string) (string, error) {

	val, ok := r.Config[key]

	if ok {
		return val, nil
	}

	return "", errors.New(key + " doesn't exist in config")

}
