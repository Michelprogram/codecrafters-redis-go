package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	// Uncomment this block to pass the first stage
	"net"
)

var data map[string]string

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	data = make(map[string]string)

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		panic(err)
	}

	defer l.Close()

	for {

		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			err := response(conn)
			if err != nil {
				log.Println(err)
				return
			}
		}()

	}

}

func createBulkString(data []byte) []byte {
	res := fmt.Sprintf("$%d\r\n%s\r\n", len(data), data)

	return []byte(res)
}

func response(conn net.Conn) error {

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

		log.Println(buffer[:size])

		commands := bytes.Split(buffer[:size], []byte("\r\n"))

		command := string(bytes.ToLower(commands[2]))

		switch command {
		case "ping":
			_, err = conn.Write([]byte("+PONG\r\n"))

		case "echo":

			_, err = conn.Write(createBulkString(commands[4]))

		case "set":
			data[string(commands[4])] = string(commands[6])
			_, err = conn.Write([]byte("+OK\r\n"))

		case "get":
			if val, ok := data[string(commands[4])]; ok {
				_, err = conn.Write(createBulkString([]byte(val)))
			} else {
				_, err = conn.Write([]byte("$-1\r\n"))
			}
		}

		if err != nil {
			return err
		}
	}

	return nil

}
