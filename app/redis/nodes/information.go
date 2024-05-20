package nodes

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Information struct {
	Role                    string
	IsMaster                bool
	MasterAddress           string
	MasterPort              int
	MasterReplicationId     string
	MasterReplicationOffset int
}

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	rand.New(src)
}

func newInformation(role string) Information {

	if role == "" {
		return Information{
			Role:                    "master",
			IsMaster:                true,
			MasterAddress:           "",
			MasterPort:              0,
			MasterReplicationId:     randomID(),
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

func randomID() string {

	var buffer bytes.Buffer

	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	for i := 0; i < 40; i++ {

		buffer.WriteByte(chars[rand.Intn(len(chars))])

	}

	return buffer.String()

}
