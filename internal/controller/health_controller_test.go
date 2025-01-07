package controller_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flohansen/auther/internal/controller"
	"github.com/stretchr/testify/assert"
)

func TestHealthController_Healthz(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {
		// given
		c := controller.NewHealthController()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/healthz", nil)

		// when
		c.Healthz(w, r)

		// then
		res := w.Result()
		body, _ := io.ReadAll(res.Body)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "OK", string(body))
	})
}
