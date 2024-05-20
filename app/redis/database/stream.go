package database

type Stream struct {
	ID    []byte
	Key   [][]byte
	Value [][]byte
}

func NewStream(id []byte, key, value []byte) Stream {

	return Stream{
		ID: id,
		Key: [][]byte{
			key,
		},
		Value: [][]byte{
			value,
		},
	}
}

func (s *Stream) Push(key, value []byte) {

	s.Key = append(s.Key, key)
	s.Value = append(s.Value, value)

}

func (s Stream) String() string {

	return ""

}
