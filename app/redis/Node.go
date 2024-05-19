package redis

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
)

type Node struct {
	*Redis
	Master net.Conn
}

func newNode(port uint, role string) *Node {
	return &Node{
		Redis:  newRedis(port, role),
		Master: nil,
	}
}

func (m *Node) ListenAndServe() error {
	l, err := net.Listen(TCP, m.Address)

	m.Listener = l

	if err != nil {
		return err
	}
	defer l.Close()

	log.Println(l.Addr())

	if err = m.handshake(); err != nil {
		return err
	}

	m.handleRequests()

	return nil
}

func (m *Node) send(data string) ([]byte, error) {

	_, err := m.Master.Write([]byte(data))
	if err != nil {
		return nil, err
	}

	received := make([]byte, 1024)
	size, err := m.Master.Read(received)
	if err != nil {
		return nil, err
	}

	return received[:size], nil

}

func (m *Node) handshake() error {

	conn, err := net.Dial(TCP, m.MasterAddress)
	if err != nil {
		return err
	}

	//defer conn.Close()

	m.Master = conn

	response, err := m.send("*1\r\n$4\r\nPING\r\n")

	if err != nil {
		return err
	}

	if !bytes.Equal(response, PONG) {
		return errors.New("Can't connected to main node at " + m.MasterAddress)
	}

	data := fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n%d\r\n", m.Port)

	response, err = m.send(data)

	if err != nil {
		return err
	}

	if !bytes.Equal(response, OK) {
		return errors.New("Can't replconf to main node at " + m.MasterAddress)
	}

	response, err = m.send("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n")

	if err != nil {
		return err
	}

	if !bytes.Equal(response, OK) {
		return errors.New("Can't replconf capa to main node at " + m.MasterAddress)
	}

	response, err = m.send("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n")

	return err
}
