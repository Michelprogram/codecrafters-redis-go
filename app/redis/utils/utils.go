package utils

import (
	"fmt"
	"net"
)

func Send(conn net.Conn, data fmt.Stringer) ([]byte, error) {

	_, err := fmt.Fprintf(conn, data.String())
	if err != nil {
		return nil, err
	}

	received := make([]byte, 1024)
	size, err := conn.Read(received)
	if err != nil {
		return nil, err
	}

	return received[:size], nil

}
