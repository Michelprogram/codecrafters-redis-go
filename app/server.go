package main

import (
	"bytes"
	"fmt"
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

		/*		for _, command := range commands {

				go response(command, conn)
			}*/

	}

}

func response(conn net.Conn) error {

	defer conn.Close()

	var err error

	buffer := make([]byte, 1024)

	size, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	command := bytes.Split(buffer[:size], []byte("\r\n"))[2]

	log.Println(string(command))

	switch string(command) {
	case "PING":
		_, err = conn.Write([]byte("+PONG\r\n"))
	}

	return err

}
