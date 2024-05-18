package main

import (
	"flag"
	"github.com/codecrafters-io/redis-starter-go/app/redis"
)

func main() {

	var port uint

	var replicaof string

	flag.UintVar(&port, "port", 6379, "Port where your server should start")
	flag.StringVar(&replicaof, "replicaof", "", "Address to the master node")

	flag.Parse()

	server := redis.NewServerFactory(port, replicaof)

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}

}
