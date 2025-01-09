package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	v1 "github.com/flohansen/auther/api/v1"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *v1.RegisterRequest) error
	LoginUser(ctx context.Context, req *v1.LoginRequest) error
}

type AuthController struct {
	userService UserService
}

func NewAuthController(userService UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req v1.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		v1.ErrorResponse(w, http.StatusBadRequest, "Request is not a valid json")
		return
	}

	if err := validateRegisterRequest(&req); err != nil {
		v1.ErrorResponse(w, http.StatusBadRequest, "Invalid username or password")
		return
	}

	if err := c.userService.RegisterUser(ctx, &req); err != nil {
		log.Printf("could not register user: %s", err)
		v1.ErrorResponse(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	v1.SuccessResponse(w, "User created")
}

func validateRegisterRequest(req *v1.RegisterRequest) error {
	if !regexp.MustCompile(`.*@.*\..*`).MatchString(req.Username) {
		return errors.New("invalid email")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(req.Password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(req.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(req.Password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(req.Password)

	isValidPassword := len(req.Password) >= 8 &&
		len(req.Password) <= 20 &&
		hasUpper &&
		hasLower &&
		hasNumber &&
		hasSpecial

	if !isValidPassword {
		return errors.New("invalid password")
	}

	return nil
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req v1.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		v1.ErrorResponse(w, http.StatusBadRequest, "Request is not a valid json")
		return
	}

	if err := c.userService.LoginUser(ctx, &req); err != nil {
		v1.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
	}
}

func (c *AuthController) Update(w http.ResponseWriter, r *http.Request) {
}
