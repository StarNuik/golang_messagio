package model

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/starnuik/golang_messagio/internal"
)

type MetricsModel struct {
	sql *pgxpool.Pool
}

type Metrics struct {
	Messages struct {
		Total      int
		LastDay    int
		LastHour   int
		LastMinute int
	}
	ProcessedTotal int
	ProcessedRatio float32
	OrphanMessages int
}

func NewMetricsModel(ctx context.Context, dbUrl string) (*MetricsModel, error) {
	pool, err := internal.NewSqlPool(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	return &MetricsModel{sql: pool}, nil
}

func (m *MetricsModel) Close() {
	m.sql.Close()
}

func queryInt(sql *pgxpool.Pool, ctx context.Context, query string, out *int) error {
	row := sql.QueryRow(ctx, query)

	err := row.Scan(out)
	return err
}

func (m *MetricsModel) Get(ctx context.Context) (Metrics, error) {
	metrics := Metrics{}

	table := []struct {
		dest  *int
		query string
	}{
		{&metrics.Messages.Total, "SELECT count(*) FROM messages;"},
		{&metrics.Messages.LastDay, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '24 hour';"},
		{&metrics.Messages.LastHour, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '1 hour';"},
		{&metrics.Messages.LastMinute, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '1 minute';"},
		{&metrics.ProcessedTotal, "SELECT count(*) FROM messages WHERE msg_is_processed=true;"},
	}

	for _, item := range table {
		err := queryInt(m.sql, ctx, item.query, item.dest)
		if err != nil {
			return metrics, err
		}
	}

	metrics.ProcessedRatio = float32(metrics.ProcessedTotal) / float32(metrics.Messages.Total)

	err := queryInt(m.sql, ctx,
		"select count(*) from messages where msg_is_processed = false and msg_created < now() - interval '1 minute';",
		&metrics.OrphanMessages)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}
