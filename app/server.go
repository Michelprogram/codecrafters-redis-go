package main

import (
	"flag"
	"github.com/codecrafters-io/redis-starter-go/app/redis/nodes"
)

func main() {

	var port uint
	var replicaof string
	var dir string
	var dbfilename string

	flag.UintVar(&port, "port", 6379, "Port where your server should start")
	flag.StringVar(&replicaof, "replicaof", "", "Address to the master node")
	flag.StringVar(&dir, "dir", "", "The path to the directory where the RDB file is stored")
	flag.StringVar(&dbfilename, "dbfilename", "", "the name of the RDB file")

	flag.Parse()

	server := nodes.NewNode(port, replicaof, dir, dbfilename)

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
	/*
		os.WriteFile("test.rdb",
			[]byte{82, 69, 68, 73, 83, 48, 48, 48, 51, 250, 9, 114, 101, 100, 105, 115, 45, 118, 101, 114, 5, 55, 46, 50, 46, 48, 250, 10, 114, 101, 100, 105, 115, 45, 98, 105, 116, 115, 192, 64, 254, 0, 251, 3, 0, 0, 9, 114, 97, 115, 112, 98, 101, 114, 114, 121, 5, 109, 97, 110, 103, 111, 0, 10, 115, 116, 114, 97, 119, 98, 101, 114, 114, 121, 5, 97, 112, 112, 108, 101, 0, 4, 112, 101, 97, 114, 6, 111, 114, 97, 110, 103, 101, 255, 170, 76, 0, 74, 110, 73, 10, 227, 10},
			0666)*/

}
