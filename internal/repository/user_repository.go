package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/flohansen/auther/internal/service"
)

type User struct {
	Username       string
	PasswordHash   []byte
	CreatedAt      time.Time
	LastModifiedAt time.Time
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

const createUserQuery = `
INSERT INTO users (username, password_hash, created_at, last_modified_at)
VALUES ($1, $2, $3, $4)
`

func (u *UserRepository) CreateUser(ctx context.Context, user service.User) error {
	_, err := u.db.ExecContext(ctx, createUserQuery, user.Username, user.PasswordHash, time.Now(), time.Now())
	return err
}
