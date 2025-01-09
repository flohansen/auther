package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

	t.Run("should return 400 BAD REQUEST if request body is invalid json", func(t *testing.T) {
		// given
		c := controller.NewAuthController(userServiceMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{
			"username": "example
		}`)))

		// when
		c.Register(w, r)

		// then
		res := w.Result()

		var m map[string]any
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&m))
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "Request is not a valid json", m["message"])
	})

	t.Run("should return 400 BAD REQUEST if username is not a valid email", func(t *testing.T) {
		// given
		c := controller.NewAuthController(userServiceMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{
			"username": "example",
			"password": "tEst!234"
		}`)))

		// when
		c.Register(w, r)

		// then
		res := w.Result()

		var m map[string]any
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&m))
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "Invalid username or password", m["message"])
	})

	t.Run("should return 400 BAD REQUEST if password don't match all rules", func(t *testing.T) {
		// given
		c := controller.NewAuthController(userServiceMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{
			"username": "example@example.com",
			"password": "password"
		}`)))

		// when
		c.Register(w, r)

		// then
		res := w.Result()

		var m map[string]any
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&m))
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "Invalid username or password", m["message"])
	})

	t.Run("should return 500 INTERNAL SERVER ERROR if user service returns error", func(t *testing.T) {
		// given
		c := controller.NewAuthController(userServiceMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{
			"username": "example@example.com",
			"password": "tEst!234"
		}`)))

		userServiceMock.EXPECT().
			RegisterUser(r.Context(), &v1.RegisterRequest{
				Username: "example@example.com",
				Password: "tEst!234",
			}).
			Return(errors.New("error"))

		// when
		c.Register(w, r)

		// then
		res := w.Result()

		var m map[string]any
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&m))
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "Could not create user", m["message"])
	})

	t.Run("should validate inputs and create user using service", func(t *testing.T) {
		// given
		c := controller.NewAuthController(userServiceMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{
			"username": "example@example.com",
			"password": "tEst!234"
		}`)))

		userServiceMock.EXPECT().
			RegisterUser(r.Context(), &v1.RegisterRequest{
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
