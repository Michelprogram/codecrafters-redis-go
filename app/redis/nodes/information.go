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
	ID                      string
	Role                    string
	IsMaster                bool
	MasterAddress           string
	MasterPort              int
	MasterReplicationId     string
	MasterReplicationOffset int
	Config                  map[string]string
}

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	rand.New(src)
}

func newInformation(role, dir, dbfilename string) Information {

	information := Information{
		ID:                      randomID(),
		MasterReplicationId:     "",
		MasterReplicationOffset: 0,
		Config: map[string]string{
			"dir":        dir,
			"dbfilename": dbfilename,
		},
	}

	if role == "" {

		information.Role = "master"
		information.IsMaster = true
		information.MasterAddress = ""
		information.MasterPort = 0
		information.MasterReplicationId = ""
		information.MasterReplicationOffset = 0

	} else {
		address := strings.Split(role, " ")

		port, _ := strconv.Atoi(address[1])

		information.Role = "slave"
		information.IsMaster = false

		information.MasterAddress = strings.Replace(role, " ", ":", 1)
		information.MasterPort = port

	}

	return information
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
