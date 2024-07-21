package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

func newMessage(req MessageReq) (Message, error) {
	if len(req.Content) <= 0 {
		return Message{}, fmt.Errorf("zero-length content")
	}

	len := min(len(req.Content), 4096)
	id, err := uuid.NewV4()
	if err != nil {
		return Message{}, fmt.Errorf("could not generate a uuid")
	}

	return Message{
		Id:      id,
		Created: time.Now().UTC(),
		Content: req.Content[:len],
	}, nil
}

func insertMessage(msg Message) error {
	tag, err := sql.Exec(
		context.Background(),
		"INSERT INTO messages (msg_id, msg_created, msg_content) VALUES ($1, $2, $3)",
		msg.Id, msg.Created, msg.Content)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		//todo: should this be a "stop the service" error or a plain log.error
		return fmt.Errorf("rowsAffected != 1")
	}
	return nil
}