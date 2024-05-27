package utils

import (
	"bytes"
	"fmt"
	"net"
	"os"
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

func ParseFile(path string) (map[string]string, error) {

	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	res := make(map[string]string)

	cursor := bytes.IndexByte(file, 251)
	end := bytes.IndexByte(file, 255)

	//First loop
	key := file[cursor+1 : end]

	keySize := int(key[3]) + 4

	valueSize := int(key[keySize]) + 1

	res[string(key[4:keySize])] = string(key[keySize : valueSize+keySize])

	cursor += valueSize + keySize

	for cursor < end-1 {
		key = file[cursor+1 : end]

		keySize = int(key[1]) + 2

		valueSize = int(key[keySize]) + 1

		res[string(key[2:keySize])] = string(key[keySize : valueSize+keySize])

		cursor += valueSize + keySize
	}

	return res, err

}
