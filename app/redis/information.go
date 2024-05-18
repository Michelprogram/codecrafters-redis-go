package redis

import (
	"fmt"
	"strings"
)

type Information struct {
	Role                    string
	IsMaster                bool
	MasterAddress           string
	MasterReplicationId     string
	MasterReplicationOffset int
}

func newInformation(role string, masterAddress string) Information {

	if role == "" {
		return Information{
			Role:                    "master",
			IsMaster:                true,
			MasterAddress:           "",
			MasterReplicationId:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
			MasterReplicationOffset: 0,
		}
	}

	return Information{
		Role:                    "slave",
		IsMaster:                false,
		MasterAddress:           strings.Replace(role, " ", ":", 1),
		MasterReplicationId:     "",
		MasterReplicationOffset: 0,
	}
}

func (i Information) String() string {
	return fmt.Sprintf("role:%smaster_replid:%smaster_repl_offset:%d", i.Role, i.MasterReplicationId, i.MasterReplicationOffset)
}
