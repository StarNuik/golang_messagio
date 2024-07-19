package internal

import (
	"context"
	"fmt"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlDriver struct {
	pool *pgxpool.Pool
}

func NewSqlDriver(pgUrl string) (*SqlDriver, error) {
	config, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		return nil, err
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &SqlDriver{
		pool: pool,
	}, nil
}

func (d *SqlDriver) InsertMessage(m Message) error {
	tag, err := d.pool.Exec(
		context.Background(),
		"INSERT INTO messages (msg_id, msg_created, msg_content, msg_processed) VALUES ($1, $2, $3, false);",
		m.Id, m.Created, m.Content)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("ERROR: InsertMessage rowsAffected != 1")
	}
	return nil
}

func (d *SqlDriver) QueryMetrics() (Metrics, error) {
	metrics := Metrics{}

	row := d.pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM messages;")

	err := row.Scan(&metrics.MessagesTotal)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (d *SqlDriver) Cleanup() {
	d.pool.Close()
}
