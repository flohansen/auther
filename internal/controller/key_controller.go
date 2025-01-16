package controller

import (
	"encoding/json"
	"log"
	"net/http"

	v1 "github.com/flohansen/auther/api/v1"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	KeyType   string `json:"kty"`
	Usage     string `json:"use"`
	KeyID     string `json:"kid"`
	Algorithm string `json:"alg"`
	Curve     string `json:"crv"`
	X         string `json:"x"`
	Y         string `json:"y"`
}

type KeyService interface {
	GetJWKS() (JWKS, error)
}

type KeyController struct {
	keyService KeyService
}

func NewKeyController(keyService KeyService) *KeyController {
	return &KeyController{
		keyService: keyService,
	}
}

func (c *KeyController) JWKS(w http.ResponseWriter, r *http.Request) {
	jwks, err := c.keyService.GetJWKS()
	if err != nil {
		log.Printf("could not get JWKS: %s", err)
		v1.ErrorResponse(w, http.StatusInternalServerError, "Could not get JWKS")
		return
	}

	json.NewEncoder(w).Encode(jwks)
}
