package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/flohansen/auther/internal/service"
)

var (
	ErrAlreadyExists = errors.New("user already exists")
	ErrNotFound      = errors.New("user not found")
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

func (r *UserRepository) CreateUser(ctx context.Context, user service.User) error {
	_, err := r.db.ExecContext(ctx, createUserQuery, user.Username, user.PasswordHash, time.Now(), time.Now())
	if err != nil {
		switch true {
		case strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_pkey"`):
			return ErrAlreadyExists
		default:
			return err
		}
	}

	return nil
}

const getUserQuery = `
SELECT username, password_hash, created_at, last_modified_at
FROM users
WHERE username = $1
LIMIT 1
`

func (r *UserRepository) GetUser(ctx context.Context, username string) (service.User, error) {
	row := r.db.QueryRowContext(ctx, getUserQuery, username)
	if err := row.Err(); err != nil {
		return service.User{}, fmt.Errorf("query error: %w", err)
	}

	var user User
	if err := row.Scan(
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastModifiedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.User{}, ErrNotFound
		}

		return service.User{}, fmt.Errorf("could not scan user row: %w", err)
	}

	return service.User{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, nil
}

const updateUserQuery = `
UPDATE users
SET password_hash = $2
WHERE username = $1
`

func (r *UserRepository) UpdateUser(ctx context.Context, user service.User) error {
	_, err := r.db.ExecContext(ctx, updateUserQuery, user.Username, user.PasswordHash)
	return err
}

const deleteUserQuery = `
DELETE FROM users
WHERE username = $1
LIMIT 1
`

func (r *UserRepository) DeleteUser(ctx context.Context, username string) error {
	_, err := r.db.ExecContext(ctx, deleteUserQuery, username)
	return err
}
