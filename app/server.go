package main

import (
	"bytes"
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
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
			err := readLoop(conn)
			if err != nil {
				panic(err)
			}
		}()
	}

}

func readLoop(conn net.Conn) error {
	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	if bytes.HasPrefix(buffer, []byte("*1\r\n$4\r\nPING\r\n")) {
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error sending: ", err.Error())
			os.Exit(1)
		}
	}

	return nil
}
