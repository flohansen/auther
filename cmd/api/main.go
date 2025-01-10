package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flohansen/auther/internal/application"
	"github.com/flohansen/auther/internal/controller"
	"github.com/flohansen/auther/internal/repository"
	"github.com/flohansen/auther/internal/service"

	_ "github.com/lib/pq"
)

var (
	configPath     = flag.String("config", "auther.config.yaml", "The path to the configuration file")
	privateKeyPath = flag.String("private-key", "private.key", "The path to the private key used for signing tokens")
)

func main() {
	flag.Parse()

	config, err := application.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("could not load config: %s", err)
	}

	db, err := sql.Open("postgres", config.Postgres.Dsn())
	if err != nil {
		log.Fatalf("could not open postgres connection: %s", err)
	}

	privateKey, err := LoadPrivateKey(*privateKeyPath)
	if err != nil {
		log.Fatalf("could not load private key: %s", err)
	}

	healthController := controller.NewHealthController()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(privateKey, userRepository)
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

func LoadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	block, _ := pem.Decode(b)
	pk, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	return pk, nil
}
