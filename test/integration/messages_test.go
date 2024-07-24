package integration_test

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMessagesGet(t *testing.T) {
	assert := assert.New(t)
	messages := model.NewMessagesModel(db)

	want := newMessage()
	_, err := db.Exec(ctx, "INSERT INTO messages (msg_id, msg_created, msg_content, msg_is_processed) VALUES ($1, $2, $3, $4);",
		want.Id, want.Created, want.Content, want.IsProcessed)
	assert.Nil(err)

	have, err := messages.Get(ctx, want.Id)
	assert.Nil(err)
	assert.Equal(want, have)
}

func TestMessagesInsert(t *testing.T) {
	assert := assert.New(t)
	messages := model.NewMessagesModel(db)

	want := newMessage()
	err := messages.Insert(ctx, want)
	assert.Nil(err)

	have, err := messages.Get(ctx, want.Id)
	assert.Nil(err)
	assert.Equal(want, have)

	// duplicate id error
	err = messages.Insert(ctx, want)
	assert.NotNil(err)
}

func TestMessagesUpdate(t *testing.T) {
	assert := assert.New(t)
	messages := model.NewMessagesModel(db)

	// insert want
	want := newMessage()
	err := messages.Insert(ctx, want)
	assert.Nil(err)

	// sanity check
	have, err := messages.Get(ctx, want.Id)
	assert.Nil(err)
	assert.Equal(want, have)

	// update is_processed
	want.IsProcessed = true
	err = messages.UpdateIsProcessed(ctx, want)
	assert.Nil(err)

	have, err = messages.Get(ctx, want.Id)
	assert.Nil(err)
	assert.Equal(want, have)
}
