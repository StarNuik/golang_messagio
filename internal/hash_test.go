package internal_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/starnuik/golang_messagio/internal"
	"github.com/stretchr/testify/assert"
)

func TestHashToString(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		from string
	}{
		{"hello world"},
		{"BIG TUNA"},
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."},
	}

	for _, tt := range tt {
		hash := sha256.Sum256([]byte(tt.from))
		want := fmt.Sprintf("%x", hash)
		have := internal.HashToString(hash)
		assert.Equal(want, have)
	}
}
