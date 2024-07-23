package internal_test

import (
	"testing"
	"time"

	"math/rand"
	"strings"

	"github.com/starnuik/golang_messagio/internal/api"
	"github.com/starnuik/golang_messagio/internal/message"
	"github.com/starnuik/golang_messagio/internal/model"
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

func timeApproxNow(in time.Time) bool {
	approxNow := time.Now().UTC().Add(-100 * time.Millisecond)
	return in.After(approxNow)
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	var req api.MessageRequest
	var msg model.Message
	var err error

	req = api.MessageRequest{}
	msg, err = message.Validate(req)
	assert.NotNil(err, "empty request must return an error")

	req = api.MessageRequest{Content: "string"}
	msg, err = message.Validate(req)
	assert.Nil(err)
	assert.Equal(req.Content, msg.Content)
	assert.True(!msg.Id.IsNil(), "message has a uuid")
	assert.True(timeApproxNow(msg.Created))

	req = api.MessageRequest{Content: randomString(1025)}
	msg, err = message.Validate(req)
	assert.NotNil(err, "max content size is 1024 bytes")
}
