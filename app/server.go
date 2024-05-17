package main

import (
	"bytes"
	"fmt"
	"log"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, 1024)

	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	if bytes.HasPrefix(buffer, []byte("*1\\r\\n$4\\r\\nPING\\r\\n")) {
		_, err = conn.Write([]byte("+PONG\\r\\n"))
		if err != nil {
			fmt.Println("Error sending: ", err.Error())
			os.Exit(1)
		}
	}

	log.Printf("read command:%s", buffer)

}
