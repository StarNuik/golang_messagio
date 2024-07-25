package integration_test

import (
	"testing"

	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMetricsGet(t *testing.T) {
	assert := assert.New(t)
	messages := model.NewMessagesModel(db)
	metrics := model.NewMetricsModel(db)

	resetDb()

	have, err := metrics.Get(ctx)
	assert.Nil(err)
	assert.Equal(0, have.Messages.Total)
	assert.Equal(0, have.ProcessedTotal)
	assert.Equal(0.0, have.ProcessedRatio)

	for range 5 {
		msg := newMessage()
		_ = messages.Insert(ctx, msg)
	}
	for range 5 {
		msg := newMessage()
		msg.IsProcessed = true
		_ = messages.Insert(ctx, msg)
	}

	have, err = metrics.Get(ctx)
	assert.Nil(err)

	assert.Equal(10, have.Messages.Total)
	assert.Equal(5, have.ProcessedTotal)
	assert.Equal(0.5, have.ProcessedRatio)
}
