package redis

import (
	"net"
)

type Main struct {
	*Redis
}

func newMain(port uint, role string) *Main {
	return &Main{
		Redis: newRedis(port, role),
	}
}

func (m *Main) ListenAndServe() error {

	l, err := net.Listen(TCP, m.Address)

	m.Listener = l

	if err != nil {
		return err
	}

	defer l.Close()

	m.handleRequests()

	return nil
}
