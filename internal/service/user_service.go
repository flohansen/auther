package service

import (
	"context"
	"errors"
	"fmt"

	v1 "github.com/flohansen/auther/api/v1"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string
	PasswordHash []byte
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req *v1.RegisterRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return fmt.Errorf("could not generate hash for password: %w", err)
	}

	return s.repo.CreateUser(ctx, User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	})
}

func (s *UserService) LoginUser(ctx context.Context, req *v1.LoginRequest) error {
	return errors.New("not implemented")
}
