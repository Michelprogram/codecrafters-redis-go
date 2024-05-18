test:
	git add . && git commit --allow-empty -m "run tests" && git push origin master

### Introduction

send-ping:
	 echo -e "*1\r\n$4\r\nPING\r\n" | netcat localhost 6379

send-echo:
	echo -e "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n" | netcat localhost 6379

send-get:
	echo -e "*2\r\n$3\r\nGET\r\n$3\r\nhey\r\n" | netcat localhost 6379

send-set:
	echo -e "*3\r\n$3\r\nSET\r\n$3\r\nhey\r\n$2\r\noi\r\n" | netcat localhost 6379

send-set-px:
	echo -e "*5\r\n$3\r\nSET\r\n$3\r\nhey\r\n$2\r\noi\r\n$2\r\npx\r\n$3\r\n100\r\n" | netcat localhost 6379

### Replication

send-info-replication:
	echo -e "*2\r\n$4\r\ninfo\r\n$11\r\nreplication\r\n" | netcat localhost $(port)

start-with-flag:
	go run app/server.go --port $(port)

start-as-master:
	go run app/server.go --port 2121

start-as-node:
	go run app/server.go --port 2122 --replicaof "localhost 2121"

send-set-to-master:
	echo -e "*3\r\n$3\r\nSET\r\n$3\r\naey\r\n$2\r\noi\r\n" | netcat localhost 2121

send-get-to-replication:
	echo -e "*2\r\n$3\r\nGET\r\n$3\r\nhey\r\n" | netcat localhost 2122

send-get-to-master:
	echo -e "*2\r\n$3\r\nGET\r\n$3\r\nhey\r\n" | netcat localhost 2121