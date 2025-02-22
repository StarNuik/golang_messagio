package model

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessagesModel struct {
	sql *pgxpool.Pool
}

type Message struct {
	Id          uuid.UUID
	Created     time.Time
	Content     string
	IsProcessed bool
}

func NewMessagesModel(pool *pgxpool.Pool) *MessagesModel {
	return &MessagesModel{sql: pool}
}

func (m *MessagesModel) Insert(ctx context.Context, msg Message) error {
	tag, err := m.sql.Exec(
		ctx,
		"INSERT INTO messages (msg_id, msg_created, msg_content, msg_is_processed) VALUES ($1, $2, $3, $4)",
		msg.Id, msg.Created, msg.Content, msg.IsProcessed)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("rows affected != 1")
	}
	return nil
}

func (m *MessagesModel) UpdateIsProcessed(ctx context.Context, msg Message) error {
	tag, err := m.sql.Exec(
		ctx,
		"UPDATE messages SET msg_is_processed=$2 WHERE msg_id=$1",
		msg.Id, msg.IsProcessed)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("rows affected != 1")
	}
	return nil
}

func (m *MessagesModel) Get(ctx context.Context, withId uuid.UUID) (Message, error) {
	row := m.sql.QueryRow(
		ctx,
		"SELECT msg_id, msg_created, msg_content, msg_is_processed FROM messages WHERE msg_id=$1",
		withId)

	var msg Message
	err := row.Scan(&msg.Id, &msg.Created, &msg.Content, &msg.IsProcessed)
	return msg, err
}
