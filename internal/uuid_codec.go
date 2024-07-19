package internal

import (
	"fmt"

	"github.com/gofrs/uuid"
)

type UuidCodec struct{}

func (*UuidCodec) Encode(value interface{}) ([]byte, error) {
	data, ok := value.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("ERROR: value is not a Uuid")
	}
	return data[:], nil
}

func (*UuidCodec) Decode(data []byte) (interface{}, error) {
	value, err := uuid.FromBytes(data)
	return value, err
}
