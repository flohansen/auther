package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	v1 "github.com/flohansen/auther/api/v1"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Time
}

type UserService interface {
	RegisterUser(ctx context.Context, req *v1.RegisterRequest) error
	Authenticate(ctx context.Context, req *v1.LoginRequest) (Tokens, error)
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

	if err := validateLoginRequest(&req); err != nil {
		log.Printf("invalid login request for user '%s': %s", req.Username, err)
		v1.ErrorResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	tokens, err := c.userService.Authenticate(ctx, &req)
	if err != nil {
		log.Printf("could not authenticate user '%s': %s", req.Username, err)
		v1.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	json.NewEncoder(w).Encode(v1.LoginResponse{
		TokenType:    "Bearer",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int(tokens.ExpiresIn.Unix()),
	})
}

func validateLoginRequest(req *v1.LoginRequest) error {
	if len(req.Username) == 0 {
		return errors.New("username is empty")
	}

	if len(req.Password) == 0 {
		return errors.New("password is empty")
	}

	return nil
}

func (c *AuthController) Update(w http.ResponseWriter, r *http.Request) {
}
