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
	"redis-cli -p $(port) replconf getack *"

send-xrange:
	redis-cli -p $(port) XADD stream_key 0-1 foo bar
	redis-cli -p $(port) XADD stream_key 0-2 bar baz
	redis-cli -p $(port) XADD stream_key 0-3 baz foo
	redis-cli -p $(port) XRANGE stream_key 0-2 0-3

send-xrange-minus:
	redis-cli -p $(port) XADD stream_key 0-1 foo bar
	redis-cli -p $(port) XADD stream_key 0-2 bar baz
	redis-cli -p $(port) XADD stream_key 0-3 baz foo
	redis-cli -p $(port) XADD stream_key 0-4 baz foo
	redis-cli -p $(port) XRANGE stream_key - 0-3

send-xrange-plus:
	redis-cli -p $(port) XADD stream_key 0-1 foo bar
	redis-cli -p $(port) XADD stream_key 0-2 bar baz
	redis-cli -p $(port) XADD stream_key 0-3 baz foo
	redis-cli -p $(port) XADD stream_key 0-4 baz foo
	redis-cli -p $(port) XRANGE stream_key 0-2 +

send-xread:
	redis-cli -p $(port) XADD stream_key 0-1 temperature 96
	redis-cli -p $(port) XREAD streams stream_key 0-0

send-xread-multiple-stream:
	redis-cli -p $(port) XADD stream_key 0-1 temperature 95
	redis-cli -p $(port) XADD other_stream_key 0-2 humidity 97
	redis-cli -p $(port) XREAD streams stream_key other_stream_key 0-0 0-1

send-xread-block:
	redis-cli -p $(port) XADD some_key 1526985054069-0 temperature 36
	redis-cli -p $(port) XREAD block 2000 streams some_key 1526985054069-0

send-xadd-block:
	redis-cli -p $(port) XADD some_key 1526985054079-0 temperature 37

## RDB

send-get-config:
	redis-cli -p $(port) config get $(key)

run-with-rdb:
	go run app/server.go --port 2121 --dir "." --dbfilename "test.rdb"