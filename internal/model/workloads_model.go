package model

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
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
	hash := internal.HashToString(load.Hash)

	tag, err := m.sql.Exec(
		ctx,
		"INSERT INTO processed_workloads (load_id, load_msg_id, load_created, load_hash) VALUES ($1, $2, $3, $4)",
		load.Id, load.MsgId, load.Created, hash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("rowsAffected != 1")
	}
	return nil
}
