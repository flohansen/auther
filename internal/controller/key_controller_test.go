package controller_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/flohansen/auther/internal/controller"
	"github.com/flohansen/auther/internal/controller/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=mocks/key_service.go -package=mocks github.com/flohansen/auther/internal/controller KeyService

func TestKeyController_JWKS(t *testing.T) {
	ctrl := gomock.NewController(t)
	keyServiceMock := mocks.NewMockKeyService(ctrl)

	t.Run("should return 500 INTERNAL SERVER ERROR if getting jwks failed", func(t *testing.T) {
		// given
		c := controller.NewKeyController(keyServiceMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		keyServiceMock.EXPECT().
			GetJWKS().
			Return(controller.JWKS{}, errors.New("error"))

		// when
		c.JWKS(w, r)

		// then
		res := w.Result()

		var m map[string]any
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&m))
		assert.Equal(t, 500, res.StatusCode)
		assert.Equal(t, "Could not get JWKS", m["message"])
	})

	t.Run("should return 200 OK and JWKS", func(t *testing.T) {
		// given
		c := controller.NewKeyController(keyServiceMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		keyServiceMock.EXPECT().
			GetJWKS().
			Return(controller.JWKS{
				Keys: []controller.JWK{
					{
						KeyType:   "Bearer",
						Usage:     "sig",
						KeyID:     "1234",
						Algorithm: "ES256",
						Curve:     "P256",
						X:         "x",
						Y:         "y",
					},
				},
			}, nil)

		// when
		c.JWKS(w, r)

		// then
		res := w.Result()

		var m map[string]any
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&m))
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, []any{
			map[string]any{
				"kty": "Bearer",
				"use": "sig",
				"kid": "1234",
				"alg": "ES256",
				"crv": "P256",
				"x":   "x",
				"y":   "y",
			},
		}, m["keys"])
	})
}
