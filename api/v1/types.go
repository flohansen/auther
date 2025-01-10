package v1

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Message{
		Success: false,
		Message: message,
	})
}

func SuccessResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Message{
		Success: true,
		Message: message,
	})
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope,omitempty"`
}

type UpdateRequest struct {
}

type UpdateResponse struct {
}
