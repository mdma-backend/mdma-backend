package types

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type (
	Hash []byte
	Salt []byte
)

type Token struct {
	Value string `json:"token"`
}

type Claims struct {
	jwt.RegisteredClaims
	RoleName    string                  `json:"role"`
	Permissions []permission.Permission `json:"permissions"`
}

type HashService interface {
	Hash(string) (Hash, Salt, error)
	HashAndCompare(string, Hash, Salt) bool
}

type TokenService interface {
	SignWithClaims(Claims) (Token, error)
	Validate(string) (*Claims, error)
}
