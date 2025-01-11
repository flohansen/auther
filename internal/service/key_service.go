package service

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"

	"github.com/flohansen/auther/internal/controller"
)

type KeyService struct {
	publicKey *ecdsa.PublicKey
}

func NewKeyService(publicKey *ecdsa.PublicKey) *KeyService {
	return &KeyService{
		publicKey: publicKey,
	}
}

func (s *KeyService) GetJWKS() (controller.JWKS, error) {
	publicKeyDER, err := x509.MarshalPKIXPublicKey(s.publicKey)
	if err != nil {
		return controller.JWKS{}, err
	}

	hash := sha256.Sum256(publicKeyDER)
	keyID := base64.RawURLEncoding.EncodeToString(hash[:8])

	jwks := controller.JWKS{
		Keys: []controller.JWK{
			{
				KeyType:   "EC",
				Usage:     "sig",
				Curve:     "P-256",
				KeyID:     keyID,
				Algorithm: "ES256",
				X:         base64.RawURLEncoding.EncodeToString(s.publicKey.X.Bytes()),
				Y:         base64.RawURLEncoding.EncodeToString(s.publicKey.Y.Bytes()),
			},
		},
	}

	return jwks, nil
}
