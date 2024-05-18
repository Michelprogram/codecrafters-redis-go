package main

import (
	"flag"
	"github.com/codecrafters-io/redis-starter-go/app/redis"
)

func main() {

	var port uint

	flag.UintVar(&port, "port", 6379, "Port where your server should start")

	flag.Parse()

	server := redis.NewServer(port)

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}

}
