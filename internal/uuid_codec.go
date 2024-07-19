package internal

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid"
)

type UuidCodec struct{}

func (*UuidCodec) Encode(value interface{}) ([]byte, error) {
	id, ok := value.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("ERROR: value is not a Uuid")
	}
	payload := id.String()
	log.Println("payload:", payload)
	return []byte(payload), nil
}

func (*UuidCodec) Decode(data []byte) (interface{}, error) {
	value, err := uuid.FromString(string(data))
	return value, err
}
