package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/flohansen/auther/internal/application"
	"github.com/flohansen/auther/internal/controller"
	"github.com/flohansen/auther/internal/repository"
	"github.com/flohansen/auther/internal/service"
)

var (
	configPath = flag.String("config", "auther.config.yaml", "The path to the configuration file")
)

func main() {
	config, err := application.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("could not load config: %s", err)
	}

	db, err := sql.Open("postgres", config.Postgres.Dsn())
	if err != nil {
		log.Fatalf("could not open postgres connection: %s", err)
	}

	healthController := controller.NewHealthController()

	userRepository := repository.NewUserRepository(db)
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
