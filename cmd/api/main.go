package main

import (
	"log"
	"net/http"

	"github.com/flohansen/auther/internal/controller"
	"github.com/flohansen/auther/internal/repository"
	"github.com/flohansen/auther/internal/service"
)

func main() {
	healthController := controller.NewHealthController()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	authController := controller.NewAuthController(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthController.Healthz)
	mux.HandleFunc("POST /api/v1/auth/register", authController.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authController.Login)
	mux.HandleFunc("PUT /api/v1/auth/update", authController.Update)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("error while listening and serving: %s", err)
	}
}
