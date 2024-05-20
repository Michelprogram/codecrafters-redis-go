test:
	git add . && git commit --allow-empty -m "run tests" && git push origin master

run-server:
	go run app/server.go --port $(port)

run-node:
	go run app/server.go --port $(port) --replicaof "localhost $(master_port)"

### Introduction

send-ping:
	redis-cli -p $(port) ping

send-echo:
	redis-cli -p $(port) echo dorian

send-get:
	redis-cli -p $(port) get $(key)

send-set:
	redis-cli -p $(port) set $(key) $(value)

send-set-px:
	redis-cli -p $(port) set $(key) $(value) px 100

### Replication

send-info-replication:
	redis-cli -p $(port) info replication

send-ack-replication:
	redis-cli -p $(port) replconf getack *