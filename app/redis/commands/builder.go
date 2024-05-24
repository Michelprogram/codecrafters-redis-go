package commands

import (
	"bytes"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/redis/database"
	"strconv"
	"strings"
)

type BuilderRESP struct {
	strings.Builder
	size int
}

func (b *BuilderRESP) updateSize() {

	b.size++

	previous := []byte(b.String())

	asciiChar := []byte(strconv.Itoa(b.size))

	previous[1] = asciiChar[0]

	b.Reset()

	b.Write(previous)

}

func (b *BuilderRESP) Bytes() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(b.String())

	return buffer.Bytes()
}

func (b *BuilderRESP) Start(resp RESP) *BuilderRESP {
	data := strings.Builder{}

	data.Write(resp)
	data.WriteByte('0')
	data.Write(CRLF)

	return &BuilderRESP{
		Builder: data,
		size:    0,
	}
}

func (b *BuilderRESP) EncodeAsArray(elements ...string) *BuilderRESP {

	b.Reset()

	b = b.Start(ARRAYS)

	for _, element := range elements {
		b.AddArgString(element)
	}

	return b
}

func (b *BuilderRESP) EncodeAsSimpleString(arg string, resp RESP) *BuilderRESP {

	data := strings.Builder{}

	data.Write(resp)
	data.WriteString(arg)
	data.Write(CRLF)

	return &BuilderRESP{
		Builder: data,
		size:    0,
	}
}

func NewBulkString[V []byte | string](arg V) *BuilderRESP {

	data := strings.Builder{}

	size := []byte(strconv.Itoa(len(arg)))

	data.Write(BULK_STRING)
	data.Write(size)
	data.Write(CRLF)
	data.Write([]byte(arg))
	data.Write(CRLF)

	return &BuilderRESP{
		Builder: data,
		size:    0,
	}
}

func (b *BuilderRESP) AddArgString(arg string) *BuilderRESP {

	defer b.Write(CRLF)

	b.updateSize()

	b.WriteByte('$')
	b.Write([]byte(strconv.Itoa(len(arg))))
	b.Write(CRLF)
	b.WriteString(arg)

	return b
}

func (b *BuilderRESP) Ping() *BuilderRESP {

	b.Reset()

	return b.Start(ARRAYS).AddArgString("PING")
}

func (b *BuilderRESP) ReplConfListeningPort(port uint) *BuilderRESP {
	b.Reset()

	portS := strconv.Itoa(int(port))

	return b.
		Start(ARRAYS).
		AddArgString("REPLCONF").
		AddArgString("listening-port").
		AddArgString(portS)
}

func (b *BuilderRESP) ReplConfCapa() *BuilderRESP {
	b.Reset()

	return b.
		Start(ARRAYS).
		AddArgString("REPLCONF").
		AddArgString("capa").
		AddArgString("psync2")
}

func (b *BuilderRESP) Psync() *BuilderRESP {

	b.Reset()

	return b.
		Start(ARRAYS).
		AddArgString("PSYNC").
		AddArgString("?").
		AddArgString("-1")
}

func (b *BuilderRESP) Ok() *BuilderRESP {

	b.Reset()

	return b.EncodeAsSimpleString("OK", SIMPLE_STRING)
}

func (b *BuilderRESP) Null() *BuilderRESP {

	b.Reset()

	return b.EncodeAsSimpleString("-1", BULK_STRING)
}

func (b *BuilderRESP) XRange(stream database.Stream) *BuilderRESP {

	b.Reset()

	b.WriteString(fmt.Sprintf("*%d", stream.Size))
	b.Write(CRLF)

	for i, id := range stream.ID {

		b.WriteString("*2")
		b.Write(CRLF)

		b.Write(NewBulkString(id.String()).Bytes())

		b.WriteString(fmt.Sprintf("*%d", stream.Size))
		b.Write(CRLF)

		b.Write(NewBulkString(stream.Key[i]).Bytes())
		b.Write(NewBulkString(stream.Value[i]).Bytes())

	}

	return b
}
