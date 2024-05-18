package redis

import (
	"fmt"
	"strconv"
	"strings"
)

type Information struct {
	Role                    string
	IsMaster                bool
	MasterAddress           string
	MasterPort              int
	MasterReplicationId     string
	MasterReplicationOffset int
}

func newInformation(role string) Information {

	if role == "" {
		return Information{
			Role:                    "master",
			IsMaster:                true,
			MasterAddress:           "",
			MasterPort:              0,
			MasterReplicationId:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
			MasterReplicationOffset: 0,
		}
	}

	address := strings.Split(role, " ")
	
	port, _ := strconv.Atoi(address[1])

	return Information{
		Role:                    "slave",
		IsMaster:                false,
		MasterAddress:           strings.Replace(role, " ", ":", 1),
		MasterPort:              port,
		MasterReplicationId:     "",
		MasterReplicationOffset: 0,
	}
}

func (i Information) String() string {
	return fmt.Sprintf("role:%smaster_replid:%smaster_repl_offset:%d", i.Role, i.MasterReplicationId, i.MasterReplicationOffset)
}
