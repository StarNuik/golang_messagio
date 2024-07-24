package model_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/peterldowns/pgmigrate"
	"github.com/starnuik/golang_messagio/internal"
	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/stretchr/testify/assert"
)

var (
	db             *pgxpool.Pool
	resetDb        func()
	testTime       time.Time
	migrationsPath = os.Getenv("TEST_MIGRATIONS_PATH")
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Panicf("Could not connect to Docker: %s", err)
	}

	poolOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=pg",
			"POSTGRES_PASSWORD=insecure",
			"POSTGRES_DB=dev",
		},
	}
	container, err := pool.RunWithOptions(&poolOpts, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	container.Expire(120)
	if err != nil {
		log.Panicf("Could not start resource: %s", err)
	}
	pgPort := container.GetPort("5432/tcp")
	pgUrl := fmt.Sprintf("postgres://pg:insecure@localhost:%s/dev", pgPort)

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

	mdb, err := sql.Open("pgx", pgUrl)
	if err != nil {
		log.Panicln(err)
	}
	defer mdb.Close()

	migrations, err := pgmigrate.Load(os.DirFS(migrationsPath))
	if err != nil {
		log.Panicln(err)
	}
	migrate := pgmigrate.NewMigrator(migrations)
	migrate.Migrate(context.Background(), mdb)

	resetDb = func() {
		_, err = mdb.Exec("DELETE FROM messages;")
		if err != nil {
			log.Panicln(err)
		}
	}

	resetDb()

	testTime = time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC)

	m.Run()
}

func TestGet(t *testing.T) {
	// resetDb()

	assert := assert.New(t)
	messages := model.NewMessagesModel(db)

	id, _ := uuid.NewV4()
	want := model.Message{
		Id:          id,
		Created:     testTime,
		Content:     "BIG TUNA",
		IsProcessed: false,
	}
	_, err := db.Exec(context.Background(), "INSERT INTO messages (msg_id, msg_created, msg_content, msg_is_processed) VALUES ($1, $2, $3, $4);",
		want.Id, want.Created, want.Content, want.IsProcessed)
	if err != nil {
		log.Panicln(err)
	}

	have, err := messages.Get(context.Background(), want.Id)
	assert.Nil(err)
	assert.Equal(want, have)
}

func TestInsert(t *testing.T) {
	// resetDb()

	assert := assert.New(t)
	messages := model.NewMessagesModel(db)

	id, _ := uuid.NewV4()
	want := model.Message{
		Id:          id,
		Created:     testTime,
		Content:     "BIG TUNA",
		IsProcessed: false,
	}
	err := messages.Insert(context.Background(), want)
	assert.Nil(err)

	have, err := messages.Get(context.Background(), want.Id)
	assert.Nil(err)
	assert.Equal(want, have)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	messages := model.NewMessagesModel(db)

	id, _ := uuid.NewV4()
	want := model.Message{
		Id:          id,
		Created:     testTime,
		Content:     "BIG TUNA",
		IsProcessed: false,
	}
	err := messages.Insert(context.Background(), want)
	assert.Nil(err)

	want.IsProcessed = true
	err = messages.Update(context.Background(), want)
	assert.Nil(err)

	have, err := messages.Get(context.Background(), want.Id)
	assert.Nil(err)
	assert.Equal(want, have)
}
