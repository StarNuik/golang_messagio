package model

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
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
	ProcessedRatio float64
	OrphanMessages int
}

func NewMetricsModel(pool *pgxpool.Pool) *MetricsModel {
	return &MetricsModel{sql: pool}
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

	metrics.ProcessedRatio = float64(metrics.ProcessedTotal) / float64(metrics.Messages.Total)
	if math.IsNaN(metrics.ProcessedRatio) {
		metrics.ProcessedRatio = 0.0
	}

	err := queryInt(m.sql, ctx,
		"select count(*) from messages where msg_is_processed = false and msg_created < now() - interval '1 minute';",
		&metrics.OrphanMessages)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}
