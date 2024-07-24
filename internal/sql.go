package internal

import (
	"context"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSqlPool(ctx context.Context, pgUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		return nil, err
	}

	config.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// same ctx?
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
