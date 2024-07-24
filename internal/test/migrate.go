package test

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/peterldowns/pgmigrate"
)

func Migrate(dbUrl string, migrationsPath string) error {
	mdb, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return err
	}
	defer mdb.Close()

	migrations, err := pgmigrate.Load(os.DirFS(migrationsPath))
	if err != nil {
		return err
	}

	migrate := pgmigrate.NewMigrator(migrations)
	_, err = migrate.Migrate(context.Background(), mdb)
	if err != nil {
		return err
	}

	return nil
}
