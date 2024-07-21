package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/starnuik/golang_messagio/internal"
)

type WorkloadsModel struct {
	sql *pgxpool.Pool
}

type Processed struct {
	Id      uuid.UUID
	MsgId   uuid.UUID
	Created time.Time
	Hash    internal.Hash256
}

func NewWorkloadsModel(pool *pgxpool.Pool) *WorkloadsModel {
	return &WorkloadsModel{sql: pool}
}

func (m *WorkloadsModel) Insert(ctx context.Context, load Processed) error {
	hex := internal.HashToString(load.Hash)

	tag, err := m.sql.Exec(
		ctx,
		"INSERT INTO processed_workloads (load_id, load_msg_id, load_created, load_hash) VALUES ($1, $2, $3, $4)",
		load.Id, load.MsgId, load.Created, hex)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("rowsAffected != 1")
	}
	return nil
}

func (m *WorkloadsModel) Get(ctx context.Context, withMsgId uuid.UUID) (Processed, error) {
	row := m.sql.QueryRow(
		ctx,
		"SELECT load_id, load_msg_id, load_created, load_hash FROM processed_workloads WHERE load_msg_id=$1",
		withMsgId)

	var load Processed
	var hex string

	err := row.Scan(&load.Id, &load.MsgId, &load.Created, &hex)
	if err != nil {
		return load, err
	}
	load.Hash, err = internal.StringToHash(hex)

	return load, err
}

func (m *WorkloadsModel) Exists(ctx context.Context, withMsgId uuid.UUID) (bool, error) {
	_, err := m.Get(ctx, withMsgId)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return false, err
}
