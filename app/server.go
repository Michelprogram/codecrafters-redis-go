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

		buffer := make([]byte, 1024)

		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		_, err = conn.Read(buffer)
		if err != nil {
			panic(err)
		}

		commands := bytes.Split(buffer, []byte("\n"))

		for _, command := range commands {

			go response(command, conn)
		}

	}

}

func response(data []byte, conn net.Conn) error {

	log.Println(string(data))

	var err error

	switch string(data) {
	case "PING\r":
		_, err = conn.Write([]byte("+PONG\r\n"))
	}

	return err

}
