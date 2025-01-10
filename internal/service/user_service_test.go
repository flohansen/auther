package service_test

import (
	"context"
	"testing"

	v1 "github.com/flohansen/auther/api/v1"
	"github.com/flohansen/auther/internal/service"
	"github.com/flohansen/auther/internal/service/mocks"
	"github.com/flohansen/auther/testhelper/matchers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=mocks/user_repsitory.go -package=mocks github.com/flohansen/auther/internal/service UserRepository

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := mocks.NewMockUserRepository(ctrl)

	t.Run("should hash the password", func(t *testing.T) {
		// given
		ctx := context.TODO()
		s := service.NewUserService(nil, userRepositoryMock)

		userRepositoryMock.EXPECT().
			CreateUser(ctx, matchers.UserMatches(
				matchers.Username("username"),
				matchers.PasswordHashOf("password"),
			)).
			Return(nil)

		// when
		err := s.RegisterUser(ctx, &v1.RegisterRequest{
			Username: "username",
			Password: "password",
		})

		// then
		assert.NoError(t, err)
	})
}
