package redis

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
)

func createBulkString(data []byte) []byte {
	res := fmt.Sprintf("$%d\r\n%s\r\n", len(data), data)

	return []byte(res)
}

type Ping struct{}

func (p Ping) Send(conn net.Conn, _ [][]byte, _ map[string]Data) error {
	_, err := conn.Write(PONG)

	return err
}

type Echo struct{}

func (e Echo) Send(conn net.Conn, args [][]byte, _ map[string]Data) error {

	_, err := conn.Write(createBulkString(args[0]))
	return err
}

type Set struct{}

func (s Set) Send(conn net.Conn, args [][]byte, database map[string]Data) error {

	key, content := string(args[0]), string(args[2])

	if len(args) > 4 {

		delay, err := strconv.Atoi(string(args[6]))

		if err != nil {
			return err
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(delay))

		database[key] = Data{
			Content: content,
			Context: ctx,
		}
	} else {
		database[key] = Data{
			Content: content,
		}
	}

	_, err := conn.Write(OK)

	return err
}

type Get struct{}

func (g Get) Send(conn net.Conn, args [][]byte, database map[string]Data) error {

	key := string(args[0])
	val, ok := database[key]

	var err error

	if ok {
		if val.Context == nil {
			_, err = conn.Write(createBulkString([]byte(val.Content)))
		} else {
			select {
			case <-val.Done():
				_, err = conn.Write(NULL)
			default:
				_, err = conn.Write(createBulkString([]byte(val.Content)))
			}
		}
	} else {
		_, err = conn.Write(NULL)
	}

	return err
}
