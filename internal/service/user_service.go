package service

import (
	"errors"

	v1 "github.com/flohansen/auther/api/v1"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) RegisterUser(req *v1.RegisterRequest) error {
	return errors.New("not implemented")
}
