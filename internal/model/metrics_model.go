package model

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricsModel struct {
	sql *pgxpool.Pool
}

type MetricsSubUnit[T any] struct {
	AllTime    T
	LastDay    T
	LastHour   T
	LastMinute T
}

type Metrics struct {
	Messages       MetricsSubUnit[int]
	Processed      MetricsSubUnit[int]
	ProcessedRatio MetricsSubUnit[float32]
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
		{&metrics.Messages.AllTime, "SELECT count(*) FROM messages;"},
		{&metrics.Messages.LastDay, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '24 hour';"},
		{&metrics.Messages.LastHour, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '1 hour';"},
		{&metrics.Messages.LastMinute, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '1 minute';"},
		{&metrics.Processed.AllTime, "SELECT count(*) FROM processed_workloads;"},
		{&metrics.Processed.LastDay, "SELECT count(*) FROM processed_workloads WHERE load_created > now() - interval '24 hour';"},
		{&metrics.Processed.LastHour, "SELECT count(*) FROM processed_workloads WHERE load_created > now() - interval '1 hour';"},
		{&metrics.Processed.LastMinute, "SELECT count(*) FROM processed_workloads WHERE load_created > now() - interval '1 minute';"},
	}

	for _, item := range table {
		err := queryInt(m.sql, ctx, item.query, item.dest)
		if err != nil {
			return metrics, err
		}
	}

	metrics.ProcessedRatio = MetricsSubUnit[float32]{
		AllTime:    float32(metrics.Processed.AllTime) / float32(metrics.Messages.AllTime),
		LastDay:    float32(metrics.Processed.LastDay) / float32(metrics.Messages.LastDay),
		LastHour:   float32(metrics.Processed.LastHour) / float32(metrics.Messages.LastHour),
		LastMinute: float32(metrics.Processed.LastMinute) / float32(metrics.Messages.LastMinute),
	}

	err := queryInt(m.sql, ctx,
		"select count(*) from messages as m where not exists (select * from processed_workloads as p where m.msg_id=p.load_msg_id) and m.msg_created < now() - interval '1 minute';",
		&metrics.OrphanMessages)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}
