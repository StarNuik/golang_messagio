package integration_test

import (
	"log"
	"testing"

	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/starnuik/golang_messagio/internal/stream"
	"github.com/stretchr/testify/require"
)

func TestDbMessageCreatedNew(t *testing.T) {
	require := require.New(t)

	log.Println("kafkaUrl", brokerUrl)
	msgCreated, err := stream.NewDbMessageCreated(brokerUrl, 10e3)
	require.Nil(err, err.Error())

	err = msgCreated.Close()
	require.Nil(err)
}

func TestDbMessageCreatedPublishAndRead(t *testing.T) {
	require := require.New(t)

	msgCreated, err := stream.NewDbMessageCreated(brokerUrl, 10e3)
	require.Nil(err)

	want := newMessage()
	readQueue := make(chan struct {
		Value model.Message
		Err   error
	})
	go func() {
		have, err := msgCreated.Read(ctx)
		readQueue <- struct {
			Value model.Message
			Err   error
		}{have, err}
	}()

	err = msgCreated.Publish(ctx, want)
	require.Nil(err, err.Error())
	have := <-readQueue
	require.Nil(have.Err)
	require.Equal(want, have.Value)

	err = msgCreated.Close()
	require.Nil(err)
}
