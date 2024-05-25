package nodes

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/redis/commands"
	"github.com/codecrafters-io/redis-starter-go/app/redis/database"
	"github.com/codecrafters-io/redis-starter-go/app/redis/utils"
	"io"
	"log"
	"net"
)

type Secondary struct {
	*Node
	Master net.Conn
}

func NewSecondary(port uint, role string) *Secondary {

	address := fmt.Sprintf("0.0.0.0:%d", port)

	return &Secondary{
		Node: &Node{
			Port:        port,
			Address:     address,
			Information: newInformation(role),
			Database:    database.NewDatabase(),
			Parser:      commands.NewParser(),
			IsPrimary:   false,
			Offset:      0,
			RDB:         "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2",
		},
	}
}

func (m *Secondary) ListenAndServe() error {

	if err := m.handshake(); err != nil {

		return err
	}

	l, err := net.Listen("tcp", m.Address)

	m.Listener = l

	if err != nil {
		return err
	}
	defer l.Close()

	go func() {
		err := m.responseFromMaster()
		if err != nil {
			log.Println(err)
		}
	}()

	m.handleRequests()

	return nil
}

func (m *Secondary) responseFromMaster() error {
	defer m.Master.Close()

	for {
		buffer := make([]byte, 1024)

		size, err := m.Master.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		log.Printf("Received : %s\n", string(buffer[:size]))

		cmd := m.ParseArgs(buffer[:size])

		log.Println("Command from master", cmd.Name)

		if cmd != nil {

			_ = cmd.Receive(m.Master, cmd.Parameters, m.Node)

			m.Offset += size

			log.Printf("Offset : %d\n", m.Offset)

		} else {
			return errors.New("ICommand " + cmd.Name + " doesn't exist")
		}

	}

	return nil
}

func (m *Secondary) handshake() error {

	var builder commands.BuilderRESP

	conn, err := net.Dial("tcp", m.MasterAddress)
	if err != nil {
		return err
	}

	m.Master = conn

	response, err := utils.Send(m.Master, builder.Ping())

	if err != nil {
		return err
	}

	if !bytes.Equal(response, []byte("+PONG\r\n")) {
		return errors.New("Can't ping to main node at " + m.MasterAddress)
	}

	response, err = utils.Send(m.Master, builder.ReplConfListeningPort(m.Port))

	if err != nil {
		return err
	}

	if !bytes.Equal(response, []byte("+OK\r\n")) {
		return errors.New("Can't replconf to main node at " + m.MasterAddress)
	}

	response, err = utils.Send(m.Master, builder.ReplConfCapa())

	if err != nil {
		return err
	}

	if !bytes.Equal(response, []byte("+OK\r\n")) {
		return errors.New("Can't replconf capa to main node at " + m.MasterAddress)
	}

	response, err = utils.Send(m.Master, builder.Psync())

	for {
		buffer := make([]byte, 1024)

		size, err := m.Master.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		log.Printf("Received redis101 : %s  %d\n", string(buffer[:size]), size)

		break
	}

	return err

}
