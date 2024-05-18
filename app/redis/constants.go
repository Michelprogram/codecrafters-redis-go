package redis

var (
	SPLITTER          = []byte("\r\n")
	OK       Response = []byte("+OK\r\n")
	PONG     Response = []byte("+PONG\r\n")
	NULL     Response = []byte("$-1\r\n")
	TCP               = "tcp"
)
