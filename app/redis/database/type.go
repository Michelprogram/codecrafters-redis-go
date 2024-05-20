package database

import (
	"context"
	"errors"
	"sync"
)

type Record struct {
	Content interface{}
	Type    string
	context.Context
}

func (r Record) String() string {

	switch r.Type {
	default:
	case "string":
		return r.Content.(string)

	case "stream":
		return r.Content.(Stream).String()

	}

	return ""
}

type Database struct {
	Data map[string]Record
	sync.Mutex
}

func NewDatabase() Database {
	return Database{
		Data: make(map[string]Record),
	}
}

func (d *Database) Add(key, value string, ctx context.Context) {

	defer d.Unlock()

	d.Lock()

	d.Data[key] = Record{
		value,
		"string",
		ctx,
	}

}

func (d *Database) AddX(id string, key, value []byte) {

	defer d.Unlock()

	d.Lock()

	stream, ok := d.Data[id]

	data := stream.Content.(Stream)

	if ok {

		data.Push(key, value)

	} else {
		d.Data[id] = Record{
			NewStream([]byte(id), key, value),
			"stream",
			nil,
		}
	}

}

func (d *Database) Get(key string) (Record, error) {

	defer d.Unlock()

	d.Lock()

	val, ok := d.Data[key]

	if ok {
		return val, nil
	}

	return Record{
		Type: "none",
	}, errors.New("Key " + key + "doesn't exist")

}
