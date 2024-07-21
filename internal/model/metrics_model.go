package model

import (
	"context"
)

type MetricsModel struct {
	model
}

type Metrics struct {
	MessagesTotal    int
	MessagesLastDay  int
	MessagesLastHour int
	// ProcessedTotal    int
	// ProcessedLastDay  int
	// ProcessedLastHout int
}

func NewMetricsModel(pgUrl string) (*MetricsModel, error) {
	model, err := newModel(pgUrl)
	return &MetricsModel{model: model}, err
}

func (m *MetricsModel) Get() (Metrics, error) {
	metrics := Metrics{}

	row := m.sql.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM messages;")

	err := row.Scan(&metrics.MessagesTotal)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}
