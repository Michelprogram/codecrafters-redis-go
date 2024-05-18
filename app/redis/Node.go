package redis

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

type Node struct {
	*Redis
	*net.TCPConn
}

func newNode(port uint, role string) *Node {
	return &Node{
		Redis:   newRedis(port, role),
		TCPConn: nil,
	}
}

func (m *Node) ListenAndServe() error {
	if err := m.handshake(); err != nil {
		return err
	}

	l, err := net.Listen(TCP, m.Address)
	if err != nil {
		return err
	}

	m.Listener = l
	m.handleRequests()

	return nil
}

func (m *Node) send(data string) ([]byte, error) {

	_, err := m.TCPConn.Write([]byte(data))
	if err != nil {
		return nil, err
	}

	received := make([]byte, 1024)
	size, err := m.TCPConn.Read(received)
	if err != nil {
		return nil, err
	}

	return received[:size], nil

}

func (m *Node) handshake() error {
	tcpServer, err := net.ResolveTCPAddr(TCP, m.MasterAddress)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP(TCP, nil, tcpServer)
	if err != nil {
		return err
	}

	m.TCPConn = conn

	response, err := m.send("*1\r\n$4\r\nPING\r\n")

	if !bytes.Equal(response, PONG) {
		return errors.New("Can't connected to main node at " + m.MasterAddress)
	}

	data := fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n%d\r\n", m.MasterPort)

	response, err = m.send(data)

	if !bytes.Equal(response, OK) {
		return errors.New("Can't replconf to main node at " + m.MasterAddress)
	}

	response, err = m.send("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n")

	if !bytes.Equal(response, OK) {
		return errors.New("Can't replconf capa to main node at " + m.MasterAddress)
	}

	return nil
}
