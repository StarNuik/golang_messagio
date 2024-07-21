package model

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricsModel struct {
	sql *pgxpool.Pool
}

type Metrics struct {
	MessagesTotal            int
	MessagesLastDay          int
	MessagesLastHour         int
	MessagesLastMinute       int
	ProcessedTotal           int
	ProcessedLastDay         int
	ProcessedLastHour        int
	ProcessedLastMinute      int
	ProcessedRatio           float32
	ProcessedRatioLastDay    float32
	ProcessedRatioLastHour   float32
	ProcessedRatioLastMinute float32
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
		{&metrics.MessagesTotal, "SELECT count(*) FROM messages;"},
		{&metrics.MessagesLastDay, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '24 hour';"},
		{&metrics.MessagesLastHour, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '1 hour';"},
		{&metrics.MessagesLastMinute, "SELECT count(*) FROM messages WHERE msg_created > now() - interval '1 minute';"},
		{&metrics.ProcessedTotal, "SELECT count(*) FROM processed_workloads;"},
		{&metrics.ProcessedLastDay, "SELECT count(*) FROM processed_workloads WHERE load_created > now() - interval '24 hour';"},
		{&metrics.ProcessedLastHour, "SELECT count(*) FROM processed_workloads WHERE load_created > now() - interval '1 hour';"},
		{&metrics.ProcessedLastMinute, "SELECT count(*) FROM processed_workloads WHERE load_created > now() - interval '1 minute';"},
	}

	for _, item := range table {
		err := queryInt(m.sql, ctx, item.query, item.dest)
		if err != nil {
			return metrics, err
		}
	}

	metrics.ProcessedRatio = float32(metrics.ProcessedTotal) / float32(metrics.MessagesTotal)
	metrics.ProcessedRatioLastDay = float32(metrics.ProcessedLastDay) / float32(metrics.MessagesLastDay)
	metrics.ProcessedRatioLastHour = float32(metrics.ProcessedLastHour) / float32(metrics.MessagesLastHour)
	metrics.ProcessedRatioLastMinute = float32(metrics.ProcessedLastMinute) / float32(metrics.MessagesLastMinute)

	return metrics, nil
}
