package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
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
	configPath = flag.String("config", "auther.config.yaml", "The path to the configuration file")
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

	privateKey, err := loadPrivateKey(config.OIDC.PrivateKey)
	if err != nil {
		log.Fatalf("could not load private key: %s", err)
	}

	publicKey, err := loadPublicKey(config.OIDC.PublicKey)
	if err != nil {
		log.Fatalf("could not load public key: %s", err)
	}

	healthController := controller.NewHealthController()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(privateKey, publicKey, userRepository)
	authController := controller.NewAuthController(userService)
	keyService := service.NewKeyService(publicKey)
	keyController := controller.NewKeyController(keyService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthController.Healthz)
	mux.HandleFunc("POST /api/v1/auth/register", authController.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authController.Login)
	mux.HandleFunc("PUT /api/v1/auth/update", authController.Update)
	mux.HandleFunc("POST /api/v1/auth/delete/{username}", authController.Delete)
	mux.HandleFunc("GET /api/v1/auth/.well-known/jwks.json", keyController.JWKS)

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("error while listening and serving: %s", err)
	}
}

func loadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
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

func loadPublicKey(path string) (*ecdsa.PublicKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	block, _ := pem.Decode(b)
	pk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	ecdsaPK, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not a ECDSA public key")
	}

	return ecdsaPK, nil
}
