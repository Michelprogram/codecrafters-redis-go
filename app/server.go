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

}
