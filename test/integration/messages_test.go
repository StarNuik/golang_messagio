package model_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest"
	"github.com/peterldowns/pgmigrate"
	"github.com/starnuik/golang_messagio/internal"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Panicf("Could not connect to Docker: %s", err)
	}

	container, err := pool.Run("postgres", "latest", []string{"POSTGRES_USER=pg", "POSTGRES_PASSWORD=insecure", "POSTGRES_DB=dev"})
	if err != nil {
		log.Panicf("Could not start resource: %s", err)
	}
	pgPort := container.GetPort("5432/tcp")
	pgUrl := fmt.Sprintf("pg:insecure@localhost:%s/dev", pgPort)

	err = pool.Retry(func() error {
		var err error
		db, err = internal.NewSqlPool(context.Background(), pgUrl)
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	})
	if err != nil {
		log.Panicf("Could not connect to database: %s", err)
	}

	defer func() {
		if err := pool.Purge(container); err != nil {
			log.Panicf("Could not purge resource: %s", err)
		}
	}()

	db, err := sql.Open("pgx", pgUrl)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	migrations, err := pgmigrate.Load(os.DirFS("../../migrations/*.sql"))
	if err != nil {
		log.Panicln(err)
	}
	migrate := pgmigrate.NewMigrator(migrations)
	migrate.Migrate(context.Background(), db)

	m.Run()
}

func TestInsert(t *testing.T) {

}
