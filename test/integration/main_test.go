package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/starnuik/golang_messagio/internal/cmd"
	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/starnuik/golang_messagio/internal/test"
)

var (
	db             *pgxpool.Pool
	resetDb        func()
	testTime       time.Time
	brokerUrl      string
	migrationsPath = os.Getenv("TEST_MIGRATIONS_PATH")
	ctx            = context.Background()
)

func newMessage() model.Message {
	id, _ := uuid.NewV4()
	return model.Message{
		Id:          id,
		Created:     testTime,
		Content:     "BIG TUNA",
		IsProcessed: false,
	}
}

func TestMain(m *testing.M) {
	docker, err := test.NewDocker()
	cmd.PanicIf(err)
	defer docker.Close()

	// kafka
	brokerUrl, err = docker.StartKafka()
	cmd.PanicIf(err)

	// postgrres
	pgUrl, err := docker.StartPostgres()
	cmd.PanicIf(err)

	db, err = docker.NewDbPool(pgUrl)
	cmd.PanicIf(err)
	defer db.Close()

	err = test.Migrate(pgUrl, migrationsPath)
	cmd.PanicIf(err)

	resetDb = func() {
		_, err = db.Exec(context.Background(), "DELETE FROM messages;")
		cmd.PanicIf(err)
	}
	resetDb()

	testTime = time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC)

	m.Run()
}
