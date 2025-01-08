package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	v1 "github.com/flohansen/auther/api/v1"
)

type UserService interface {
	RegisterUser(req *v1.RegisterRequest) error
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
	var req v1.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		v1.ErrorResponse(w, http.StatusBadRequest, "Request is not a valid json")
		return
	}

	if err := validateRegisterRequest(&req); err != nil {
		v1.ErrorResponse(w, http.StatusBadRequest, "Invalid username or password")
		return
	}

	if err := c.userService.RegisterUser(&req); err != nil {
		v1.ErrorResponse(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	v1.Response(w, "User created")
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
}

func (c *AuthController) Update(w http.ResponseWriter, r *http.Request) {
}
