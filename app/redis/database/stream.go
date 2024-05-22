package database

import (
	"errors"
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

func NewStream() *Stream {

	return &Stream{
		ID:    make([]ID, 0),
		Key:   make([][]byte, 0),
		Value: make([][]byte, 0),
		Size:  0,
	}
}

func (s *Stream) Push(id, key, value []byte) error {

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

	var err error

	var sequenceNumber int

	lastElement := s.ID[s.Size]

	stringID := string(id)

	infoIds := strings.Split(stringID, "-")

	millisecondsTime, err := strconv.Atoi(infoIds[0])

	if err != nil {
		return 0, 0, err
	}

	if infoIds[1] == "*" {

		if s.Size == 0 {
			sequenceNumber = 0
		} else {
			sequenceNumber = lastElement.SequenceNumber + 1
		}

	} else {
		sequenceNumber, err = strconv.Atoi(infoIds[1])

		if err != nil {
			return 0, 0, err
		}
	}

	if stringID == "0-0" {
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
