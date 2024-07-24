package integration_test

import (
	"testing"

	"github.com/starnuik/golang_messagio/internal/stream"
	"github.com/stretchr/testify/assert"
)

func TestNewDbMessageCreated(t *testing.T) {
	assert := assert.New(t)

	_, err := stream.NewDbMessageCreated(brokerUrl, 10e3)
	assert.Nil(err)
}
