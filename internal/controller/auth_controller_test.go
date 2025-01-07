package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/flohansen/auther/api/v1"
	"github.com/flohansen/auther/internal/controller"
	"github.com/flohansen/auther/internal/controller/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=mocks/user_service.go -package=mocks github.com/flohansen/auther/internal/controller UserService

func TestAuthController_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	userServiceMock := mocks.NewMockUserService(ctrl)

	t.Run("should validate inputs and create user using service", func(t *testing.T) {
		// given
		c := controller.NewAuthController(userServiceMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{
			"username": "example@example.com",
			"password": "tEst!234"
		}`)))

		userServiceMock.EXPECT().
			RegisterUser(&v1.RegisterRequest{
				Username: "example@example.com",
				Password: "tEst!234",
			}).
			Return(nil)

		// when
		c.Register(w, r)

		// then
		res := w.Result()

		var m map[string]any
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&m))
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "User created", m["message"])
	})
}

func TestAuthController_Login(t *testing.T) {
}

func TestAuthController_Update(t *testing.T) {
}
