package database

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

type ID struct {
	Id               []byte
	MillisecondsTime int
	SequenceNumber   int
}

func NewId(id []byte, millisecondsTime, sequenceNumber int) ID {

	return ID{
		id,
		millisecondsTime,
		sequenceNumber,
	}

}

type Stream struct {
	ID    []ID
	Key   [][]byte
	Value [][]byte
	Size  int
}

func NewStream(id ID, key, value []byte) *Stream {

	return &Stream{
		ID: []ID{
			id,
		},
		Key: [][]byte{
			key,
		},
		Value: [][]byte{
			value,
		},
		Size: 0,
	}
}

func (s *Stream) Push(id, key, value []byte) error {

	log.Printf("Try to push as id %s as key %s as value %s", string(id), string(key), string(value))

	ms, sn, err := s.CouldInsert(id)

	if err != nil {
		return err
	}

	s.Key = append(s.Key, key)
	s.Value = append(s.Value, value)

	s.ID = append(s.ID, NewId(id, ms, sn))

	s.Size++

	return nil

}

func (s *Stream) CouldInsert(id []byte) (int, int, error) {

	stringID := string(id)

	infoIds := strings.Split(stringID, "-")

	millisecondsTime, err := strconv.Atoi(infoIds[0])

	if err != nil {
		return 0, 0, err
	}

	sequenceNumber, err := strconv.Atoi(infoIds[1])

	if err != nil {
		return 0, 0, err
	}

	log.Println("Size : ", s.Size)

	lastElement := s.ID[s.Size]

	log.Printf("Last element as id %s as ms %d as sn %d", string(lastElement.Id), lastElement.MillisecondsTime, lastElement.SequenceNumber)

	if millisecondsTime == 0 && sequenceNumber == 0 {
		return 0, 0, errors.New("-ERR The ID specified in XADD must be greater than 0-0\r\n")
	}

	if millisecondsTime == lastElement.MillisecondsTime && sequenceNumber <= lastElement.SequenceNumber {
		return 0, 0, errors.New("-ERR The ID specified in XADD is equal or smaller than the target stream top item\r\n")
	}

	if millisecondsTime < lastElement.MillisecondsTime {
		return 0, 0, errors.New("-ERR The ID specified in XADD is equal or smaller than the target stream top item\r\n")
	}

	return millisecondsTime, sequenceNumber, nil

}

func (s Stream) String() string {

	return ""

}
