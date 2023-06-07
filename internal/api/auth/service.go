package auth

import (
	"bytes"
	"crypto/rand"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

type JWTService struct {
	SigningMethod jwt.SigningMethod
	Secret        []byte
	Leeway        time.Duration
}

func (s JWTService) SignWithClaims(claims Claims) (string, error) {
	token := jwt.NewWithClaims(s.SigningMethod, claims)
	return token.SignedString(s.Secret)
}

func (s JWTService) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	}, jwt.WithLeeway(s.Leeway))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

type Argon2IDService struct {
	SaltLen uint32
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

func (s Argon2IDService) Hash(password string) (Hash, Salt, error) {
	salt := make([]byte, s.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	hash := argon2.IDKey([]byte(password), salt, s.Time, s.Memory, s.Threads, s.KeyLen)
	return Hash(hash), Salt(salt), nil
}

func (s Argon2IDService) HashAndCompare(password string, hash Hash, salt Salt) bool {
	pwHash := argon2.IDKey([]byte(password), salt, s.Time, s.Memory, s.Threads, s.KeyLen)
	return bytes.Equal(hash, pwHash)
}
