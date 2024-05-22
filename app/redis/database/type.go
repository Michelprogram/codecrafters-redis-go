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
		return r.Content.(*Stream).String()

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

func (d *Database) AddX(key, id string, Xkey, Xvalue []byte) error {

	defer d.Unlock()

	d.Lock()

	stream, ok := d.Data[key]

	if ok {
		data := stream.Content.(*Stream)
		err := data.Push([]byte(id), Xkey, Xvalue)
		if err != nil {
			return err
		}

	} else {

		streamInit := NewStream()

		err := streamInit.Push([]byte(id), Xkey, Xvalue)

		if err != nil {
			return err
		}

		d.Data[key] = Record{
			streamInit,
			"stream",
			nil,
		}
	}

	return nil

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
