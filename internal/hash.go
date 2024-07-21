package internal

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hash256 = [sha256.Size]byte

func HashToString(hash Hash256) string {
	return hex.EncodeToString(hash[:])
}

func NewHash(data string) Hash256 {
	return sha256.Sum256([]byte(data))
}
