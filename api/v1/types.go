package v1

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorMessage{
		Success: false,
		Message: message,
	})
}

func Response(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ErrorMessage{
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
}

type LoginResponse struct {
}

type UpdateRequest struct {
}

type UpdateResponse struct {
}
