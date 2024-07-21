package model

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/starnuik/golang_messagio/internal"
)

type model struct {
	sql *pgxpool.Pool
}

func newModel(pgUrl string) (model, error) {
	sql, err := internal.NewSqlPool(pgUrl)
	return model{sql: sql}, err
}

func (m *model) Close() {
	m.sql.Close()
}
