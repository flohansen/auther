package controller

import "net/http"

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
}

func (c *AuthController) Update(w http.ResponseWriter, r *http.Request) {
}
