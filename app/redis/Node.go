package redis

import (
	"bytes"
	"errors"
	"net"
)

type Node struct {
	*Redis
}

func newNode(port uint, role string) *Node {
	return &Node{
		newRedis(port, role),
	}
}

func (m *Node) ListenAndServe() error {

	err := m.handshake()

	if err != nil {
		return err
	}

	l, err := net.Listen(TCP, m.Address)

	m.Listener = l

	if err != nil {
		return err
	}

	defer l.Close()

	m.handleRequests()

	return nil
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
	defer conn.Close()

	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		return err
	}

	received := make([]byte, 1024)
	size, err := conn.Read(received)
	if err != nil {
		return err
	}

	if bytes.Equal(received[:size], PONG) {
		return nil
	}

	return errors.New("Can't connected to main node at " + m.MasterAddress)

}
