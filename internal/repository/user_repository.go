package repository

import (
	"errors"

	"github.com/flohansen/auther/internal/service"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) CreateUser(user service.User) error {
	return errors.New("not implemented")
}
