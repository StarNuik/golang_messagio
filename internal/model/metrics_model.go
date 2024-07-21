package model

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricsModel struct {
	sql *pgxpool.Pool
}

type Metrics struct {
	MessagesTotal    int
	MessagesLastDay  int
	MessagesLastHour int
	// ProcessedTotal    int
	// ProcessedLastDay  int
	// ProcessedLastHout int
}

func NewMetricsModel(pool *pgxpool.Pool) *MetricsModel {
	return &MetricsModel{sql: pool}
}

func (m *MetricsModel) Get(ctx context.Context) (Metrics, error) {
	metrics := Metrics{}

	row := m.sql.QueryRow(ctx,
		"SELECT COUNT(*) FROM messages;")

	err := row.Scan(&metrics.MessagesTotal)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}
