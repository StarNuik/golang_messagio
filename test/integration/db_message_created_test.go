package integration_test

import (
	"testing"

	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/starnuik/golang_messagio/internal/stream"
	"github.com/stretchr/testify/require"
)

func TestDbMessageCreatedNew(t *testing.T) {
	require := require.New(t)

	msgCreated, err := stream.NewDbMessageCreated(brokerUrl, 10e3)
	require.Nil(err)

	err = msgCreated.Close()
	require.Nil(err)
}

func TestDbMessageCreatedPublishAndRead(t *testing.T) {
	require := require.New(t)

	msgCreated, err := stream.NewDbMessageCreated(brokerUrl, 10e3)
	require.Nil(err)

	want := newMessage()
	readQueue := make(chan model.Message)
	go func() {
		have, err := msgCreated.Read(ctx)
		require.Nil(err)
		readQueue <- have
	}()

	err = msgCreated.Publish(ctx, want)
	require.Nil(err)
	have := <-readQueue
	require.Equal(want, have)

	err = msgCreated.Close()
	require.Nil(err)
}
