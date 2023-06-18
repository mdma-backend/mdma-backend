package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type (
	Hash []byte
	Salt []byte
)

type Token struct {
	Value string `json:"token"`
}

type AccountType string

const (
	UserAccountType    AccountType = "user"
	ServiceAccountType AccountType = "service"
)

type Claims struct {
	jwt.RegisteredClaims
	AccountType AccountType `json:"accountType"`
	AccountID   uint        `json:"accountID"`
}

type HashService interface {
	Hash(string) (Hash, Salt, error)
	HashAndCompare(string, Hash, Salt) bool
}

type TokenService interface {
	SignWithClaims(Claims) (Token, error)
	Validate(string) (*Claims, error)
}
