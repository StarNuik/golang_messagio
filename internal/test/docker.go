package test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/starnuik/golang_messagio/internal"
)

const pgUser = "pg"
const pgPass = "insecure"
const pgDb = "test"

type Docker struct {
	pool   *dockertest.Pool
	cts    []*dockertest.Resource
	expire uint
}

func NewDocker(expireSeconds uint) (*Docker, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("Could not construct pool: %w", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("Could not connect to Docker: %w", err)
	}

	pool.MaxWait = time.Duration(expireSeconds) * time.Second

	return &Docker{
		pool:   pool,
		expire: expireSeconds,
	}, nil
}

func (d *Docker) Close() {
	pool := d.pool
	for _, container := range d.cts {
		if err := pool.Purge(container); err != nil {
			log.Panicf("Could not purge resource: %s", err)
		}
	}
}

func (d *Docker) Retry(op func() error) error {
	return d.pool.Retry(op)
}

func (d *Docker) NewDbPool(pgUrl string) (*pgxpool.Pool, error) {
	pool := d.pool

	var dbPool *pgxpool.Pool
	err := pool.Retry(func() error {
		var err error
		dbPool, err = internal.NewSqlPool(context.Background(), pgUrl)
		if err != nil {
			return err
		}
		return dbPool.Ping(context.Background())
	})
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %w", err)
	}

	return dbPool, nil
}

func (d *Docker) StartPostgres() (string, error) {
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Name:       "dockertest-postgres",
		Env: []string{
			"POSTGRES_USER=" + pgUser,
			"POSTGRES_PASSWORD=" + pgPass,
			"POSTGRES_DB=" + pgDb,
		},
	}

	container, err := d.startContainer(opts)
	if err != nil {
		return "", err
	}

	pgPort := container.GetPort("5432/tcp")
	pgUrl := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", pgUser, pgPass, pgPort, pgDb)

	return pgUrl, nil
}

func (d *Docker) StartKafka() (string, error) {
	opts := dockertest.RunOptions{
		Repository: "bitnami/kafka",
		Tag:        "3.5",
		Name:       "dockertest-kafka",
		Env: []string{
			"KAFKA_CFG_NODE_ID=0",
			"KAFKA_CFG_PROCESS_ROLES=controller,broker",
			"KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093",
			"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
			"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@localhost:9093",
			"KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER",
			"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true",
		},
	}
	container, err := d.startContainer(opts)
	if err != nil {
		return "", err
	}

	port := container.GetPort("9092/tcp")
	url := fmt.Sprintf("localhost:%s", port)
	return url, nil
}

func (d *Docker) startContainer(runOpts dockertest.RunOptions) (*dockertest.Resource, error) {
	pool := d.pool
	hostConfig := func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	}

	container, err := pool.RunWithOptions(&runOpts, hostConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not start container: %w", err)
	}

	container.Expire(d.expire)
	d.cts = append(d.cts, container)

	return container, nil
}
