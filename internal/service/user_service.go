package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	v1 "github.com/flohansen/auther/api/v1"
	"github.com/flohansen/auther/internal/controller"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	Username     string
	PasswordHash []byte
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, username string) (User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, username string) error
}

type UserService struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	repo       UserRepository
}

func NewUserService(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, repo UserRepository) *UserService {
	return &UserService{
		privateKey: privateKey,
		publicKey:  publicKey,
		repo:       repo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, req *v1.RegisterRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return fmt.Errorf("could not generate hash for password: %w", err)
	}

	return s.repo.CreateUser(ctx, User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	})
}

type Claims struct {
	jwt.RegisteredClaims
}

func (s *UserService) Authenticate(ctx context.Context, req *v1.LoginRequest) (controller.Tokens, error) {
	user, err := s.repo.GetUser(ctx, req.Username)
	if err != nil {
		return controller.Tokens{}, fmt.Errorf("could not get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password)); err != nil {
		return controller.Tokens{}, ErrInvalidCredentials
	}

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   req.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	accessTokenString, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		return controller.Tokens{}, fmt.Errorf("could not sign token: %w", err)
	}

	refreshToken := make([]byte, 32)
	if _, err := rand.Read(refreshToken); err != nil {
		return controller.Tokens{}, fmt.Errorf("could not generate refresh token: %w", err)
	}

	refreshTokenString := hex.EncodeToString(refreshToken)

	return controller.Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    time.Now().Add(24 * 7 * time.Hour),
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, tokenString string, req *v1.UpdateRequest) error {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.publicKey, nil
	})
	if err != nil {
		return fmt.Errorf("could not parse token: %w", err)
	}

	if req.Username != claims.Subject {
		return errors.New("unauthorized")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return fmt.Errorf("could not generate hash for password: %w", err)
	}

	return s.repo.UpdateUser(ctx, User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	})
}

func (s *UserService) DeleteUser(ctx context.Context, tokenString string, username string) error {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.publicKey, nil
	})
	if err != nil {
		return fmt.Errorf("could not parse token: %w", err)
	}

	if username != claims.Subject {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteUser(ctx, username)
}
