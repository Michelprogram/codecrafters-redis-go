package database

import (
	"context"
	"errors"
	"sync"
)

type data struct {
	Content string
	context.Context
}

type Database struct {
	Data map[string]data
	sync.Mutex
}

func NewDatabase() Database {
	return Database{
		Data: make(map[string]data),
	}
}

func (d *Database) Add(key, value string, ctx context.Context) {

	defer d.Unlock()

	d.Lock()

	d.Data[key] = data{
		value,
		ctx,
	}

}

func (d *Database) Get(key string) (data, error) {

	defer d.Unlock()

	d.Lock()

	val, ok := d.Data[key]

	if ok {
		return val, nil
	}

	return data{}, errors.New("Key " + key + "doesn't exist")

}
