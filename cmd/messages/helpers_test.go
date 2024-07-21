package main

import (
	"testing"
	"time"

	"math/rand"
	"strings"

	"github.com/stretchr/testify/assert"
)

// https://gosamples.dev/random-string/
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)

	var req MessageReq
	var msg Message
	var err error

	req = MessageReq{}
	msg, err = newMessage(req)
	assert.NotNil(err, "empty request must return an error")

	req = MessageReq{"string"}
	msg, err = newMessage(req)
	assert.Nil(err)
	assert.Equal(req.Content, msg.Content)
	assert.True(!msg.Id.IsNil(), "message has a uuid")
	assert.True(msg.Created.After(time.Now().UTC().Add(-100 * time.Millisecond)))

	req = MessageReq{randomString(8192)}
	msg, err = newMessage(req)
	assert.Nil(err)
	assert.True(len(msg.Content) == 4096, "trim message length")
}
