package containers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	pgUsername = "test"
	pgPassword = "test"
	pgDatabase = "test"
)

type StartedPostgresContainer struct {
	t         *testing.T
	container testcontainers.Container
}

func (c StartedPostgresContainer) Dsn() string {
	host, err := c.container.Host(context.Background())
	if err != nil {
		c.t.Fatalf("could not get container host: %s", err)
	}

	port, err := c.container.MappedPort(context.Background(), nat.Port("5432"))
	if err != nil {
		c.t.Fatalf("could not get mapped container port: %s", err)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pgUsername, pgPassword, host, port.Int(), pgDatabase)
}

func StartPostgres(t *testing.T, opts ...StartPostgresOpt) *StartedPostgresContainer {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:15",
			Env: map[string]string{
				"POSTGRES_USER":     pgUsername,
				"POSTGRES_PASSWORD": pgPassword,
				"POSTGRES_DB":       pgDatabase,
			},
			HostConfigModifier: func(cfg *container.HostConfig) {
				cfg.AutoRemove = true
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(time.Minute),
			),
		},
	})
	if err != nil {
		t.Fatalf("could not start test container: %s", err)
	}

	t.Cleanup(func() {
		container.Terminate(context.Background())
	})

	postgresContainer := &StartedPostgresContainer{
		t:         t,
		container: container,
	}

	db, err := sql.Open("postgres", postgresContainer.Dsn())
	if err != nil {
		t.Fatalf("could not open postgres connection: %s", err)
	}

	for _, opt := range opts {
		opt(t, db)
	}

	return postgresContainer
}

type StartPostgresOpt func(*testing.T, *sql.DB)

func WithMigrations(dir string) StartPostgresOpt {
	return func(t *testing.T, db *sql.DB) {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			t.Fatalf("could not create postgres driver: %s", err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", dir),
			"postgres",
			driver,
		)
		if err != nil {
			t.Fatalf("could not create migrate instance: %s", err)
		}

		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			t.Fatalf("could not migrate: %s", err)
		}
	}
}
