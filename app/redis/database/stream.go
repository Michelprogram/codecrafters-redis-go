package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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

func (i ID) String() string {

	return fmt.Sprintf("%d-%d", i.MillisecondsTime, i.SequenceNumber)

}

type Stream struct {
	ID    []ID
	Key   [][]byte
	Value [][]byte
	Size  int
	link  map[int]int
}

func NewStream() *Stream {

	return &Stream{
		ID:    make([]ID, 0),
		Key:   make([][]byte, 0),
		Value: make([][]byte, 0),
		Size:  0,
		link:  map[int]int{},
	}
}

func (s *Stream) Push(id, key, value []byte) (*ID, error) {

	ms, sn, err := s.CouldInsert(id)

	if err != nil {
		return nil, err
	}

	s.Key = append(s.Key, key)
	s.Value = append(s.Value, value)

	res := NewId(id, ms, sn)

	s.ID = append(s.ID, res)

	s.link[ms]++

	s.Size++

	return &res, nil

}

func (s *Stream) CouldInsert(id []byte) (int, int, error) {

	var err error

	var sequenceNumber int

	stringID := string(id)

	if stringID == "*" {
		buffer := []byte(strconv.FormatInt(time.Now().UnixMilli(), 10) + "-0")
		return s.CouldInsert(buffer)
	}

	if stringID == "0-0" {
		return 0, 0, errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}

	infoIds := strings.Split(stringID, "-")

	millisecondsTime, err := strconv.Atoi(infoIds[0])

	if err != nil {
		return 0, 0, err
	}

	if infoIds[1] != "*" {

		sequenceNumber, err = strconv.Atoi(infoIds[1])

		if err != nil {
			return 0, 0, err
		}

	} else if data, ok := s.link[millisecondsTime]; ok {
		sequenceNumber = data
	} else {
		if millisecondsTime == 0 {
			s.link[millisecondsTime] = 1
		} else {
			s.link[millisecondsTime] = sequenceNumber
		}

		sequenceNumber = s.link[millisecondsTime]

	}

	if s.Size > 0 {
		lastElement := s.ID[s.Size-1]

		if millisecondsTime == lastElement.MillisecondsTime && sequenceNumber <= lastElement.SequenceNumber {
			return 0, 0, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		}

		if millisecondsTime < lastElement.MillisecondsTime {
			return 0, 0, errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		}
	}

	return millisecondsTime, sequenceNumber, nil

}

func (s Stream) String() string {

	return ""

}

func (s Stream) Range(start, end []byte) (*Stream, error) {

	startInfo := strings.Split(string(start), "-")

	endInfo := strings.Split(string(end), "-")

	startMS, err := strconv.Atoi(startInfo[0])

	if err != nil {
		return nil, err
	}

	startSN, err := strconv.Atoi(startInfo[1])

	if err != nil {
		return nil, err
	}

	endMS, err := strconv.Atoi(endInfo[0])

	if err != nil {
		return nil, err
	}

	endSN, err := strconv.Atoi(endInfo[1])

	if err != nil {
		return nil, err
	}

	stream := NewStream()

	for i := 0; i < s.Size; i++ {

		element := s.ID[i]

		if element.MillisecondsTime >= startMS && element.MillisecondsTime <= endMS && element.SequenceNumber >= startSN && element.SequenceNumber <= endSN {
			stream.ID = append(stream.ID, element)
			stream.Key = append(stream.Key, s.Key[i])
			stream.Value = append(stream.Value, s.Value[i])
			stream.Size++
		}

	}

	return stream, nil

}

func (s Stream) RangeFromBeginning(end []byte) (*Stream, error) {

	endInfo := strings.Split(string(end), "-")

	endMS, err := strconv.Atoi(endInfo[0])

	if err != nil {
		return nil, err
	}

	endSN, err := strconv.Atoi(endInfo[1])

	if err != nil {
		return nil, err
	}

	stream := NewStream()

	element := s.ID[stream.Size]

	for element.MillisecondsTime >= endMS && element.SequenceNumber >= endSN {

		stream.ID = append(stream.ID, element)
		stream.Key = append(stream.Key, s.Key[stream.Size])
		stream.Value = append(stream.Value, s.Value[stream.Size])
		stream.Size++

		element = s.ID[stream.Size]

	}

	return stream, nil

}
