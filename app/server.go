package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	// Uncomment this block to pass the first stage
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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
				panic(err)
			}
		}()

	}

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

		commands := bytes.Split(buffer[:size], []byte("\r\n"))

		command := string(bytes.ToLower(commands[2]))

		log.Println(command)

		switch command {
		case "ping":
			_, err = conn.Write([]byte("+PONG\r\n"))

		case "echo":

			res := fmt.Sprintf("$%d\r\n%s\r\n", len(commands[4]), commands[4])

			_, err = conn.Write([]byte(res))

		}

		if err != nil {
			return err
		}
	}

	return nil

}
