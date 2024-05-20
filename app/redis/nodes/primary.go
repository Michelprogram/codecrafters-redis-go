package nodes

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/redis/commands"
	"github.com/codecrafters-io/redis-starter-go/app/redis/database"
	"net"
)

type Primary struct {
	*Node
}

func NewPrimary(port uint, role string) *Primary {

	address := fmt.Sprintf("0.0.0.0:%d", port)

	return &Primary{
		Node: &Node{
			Port:         port,
			Address:      address,
			Information:  newInformation(role),
			Database:     database.NewDatabase(),
			Replications: make([]net.Conn, 0, 10),
			Parser:       commands.NewParser(),
			IsPrimary:    true,
			Offset:       0,
			RDB:          "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2",
		},
	}
}

func (m *Primary) ListenAndServe() error {

	l, err := net.Listen("tcp", m.Address)

	m.Listener = l

	if err != nil {
		return err
	}

	defer l.Close()

	m.handleRequests()

	return nil
}
