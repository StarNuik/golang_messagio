package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlDriver struct {
	pool *pgxpool.Pool
}

func NewSqlDriver(pgUrl string) (*SqlDriver, error) {
	pool, err := pgxpool.New(context.Background(), pgUrl)
	if err != nil {
		return nil, fmt.Errorf("ERROR: could not connect to postgres: %v", err)
	}

	return &SqlDriver{
		pool: pool,
	}, nil
}

func (d *SqlDriver) Cleanup() {
	d.pool.Close()
}
