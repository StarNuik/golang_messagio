package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Hash256 = [sha256.Size]byte

func HashToString(hash Hash256) string {
	return hex.EncodeToString(hash[:])
}

func NewHash(data string) Hash256 {
	return sha256.Sum256([]byte(data))
}

func StringToHash(from string) (Hash256, error) {
	hash := Hash256{}
	bytes, err := hex.DecodeString(from)
	if err != nil {
		return hash, nil
	}
	if len(bytes) != 32 {
		return hash, fmt.Errorf("string to hash: len != 32")
	}
	copy(hash[:], bytes)
	return hash, nil
}
