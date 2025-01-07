package controller

import "net/http"

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
