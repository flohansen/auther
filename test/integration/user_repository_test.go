package integration_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/flohansen/auther/internal/repository"
	"github.com/flohansen/auther/internal/service"
	"github.com/flohansen/auther/testhelper/containers"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestUserRepository_CreateUser(t *testing.T) {
	postgres := containers.StartPostgres(t, containers.WithMigrations("../../sql/migrations"))

	db, err := sql.Open("postgres", postgres.Dsn())
	if err != nil {
		t.Fatalf("could not open postgres connection: %s", err)
	}

	setup := createSetup(t, db)
	createUser := createCreateUser(t, db)

	t.Run("should insert user", func(t *testing.T) {
		setup()

		// given
		r := repository.NewUserRepository(db)

		// when
		err := r.CreateUser(context.TODO(), service.User{
			Username:     "username",
			PasswordHash: []byte("password_hash"),
		})

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error if user already exists", func(t *testing.T) {
		setup()

		// given
		r := repository.NewUserRepository(db)

		createUser(repository.User{
			Username:       "username",
			PasswordHash:   []byte("password_hash"),
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		})

		// when
		err := r.CreateUser(context.TODO(), service.User{
			Username:     "username",
			PasswordHash: []byte("password_hash"),
		})

		// then
		assert.ErrorContains(t, err, "user already exists")
	})
}

func createSetup(t *testing.T, db *sql.DB) func() {
	return func() {
		if _, err := db.Exec("truncate users"); err != nil {
			t.Fatalf("could not cleanup users table: %s", err)
		}
	}
}

func createCreateUser(t *testing.T, db *sql.DB) func(repository.User) {
	return func(user repository.User) {
		_, err := db.Exec(`
			insert into users (username, password_hash, created_at, last_modified_at)
			values ($1, $2, $3, $4)
		`, user.Username, user.PasswordHash, user.CreatedAt, user.LastModifiedAt)
		if err != nil {
			t.Fatalf("could not insert user: %s", err)
		}
	}
}
