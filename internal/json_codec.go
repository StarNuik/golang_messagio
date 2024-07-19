package internal

import "encoding/json"

// todo: is this the correct way? it looks v hacky
type JsonCodec[T any] struct{}

func (_ *JsonCodec[T]) Encode(value interface{}) ([]byte, error) {
	data, err := json.Marshal(value)
	return data, err
}

func (_ *JsonCodec[T]) Decode(data []byte) (interface{}, error) {
	var t T
	err := json.Unmarshal(data, &t)
	return t, err
}
