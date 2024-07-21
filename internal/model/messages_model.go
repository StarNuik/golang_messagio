package model

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type MessagesModel struct {
	model
}

type Message struct {
	Id        uuid.UUID
	Created   time.Time
	Content   string
	Processed bool
}

func NewMessagesModel(pgUrl string) (*MessagesModel, error) {
	model, err := newModel(pgUrl)
	return &MessagesModel{model: model}, err
}

func (m *MessagesModel) Insert(msg Message) error {
	tag, err := m.sql.Exec(
		context.Background(),
		"INSERT INTO messages (msg_id, msg_created, msg_content) VALUES ($1, $2, $3)",
		msg.Id, msg.Created, msg.Content)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("rowsAffected != 1")
	}
	return nil
}
