package redis

import "fmt"

type Information struct {
	Role                    string
	MasterReplicationId     string
	MasterReplicationOffset int
}

func newInformation(role string) Information {

	if role == "master" {
		return Information{
			Role:                    role,
			MasterReplicationId:     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
			MasterReplicationOffset: 0,
		}
	}

	return Information{
		Role:                    "slave",
		MasterReplicationId:     "",
		MasterReplicationOffset: 0,
	}
}

func (i Information) String() string {
	return fmt.Sprintf("role:%smaster_replid:%smaster_repl_offset:%d", i.Role, i.MasterReplicationId, i.MasterReplicationOffset)
}
