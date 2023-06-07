package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mdma-backend/mdma-backend/internal/api/role"
	"github.com/mdma-backend/mdma-backend/internal/api/role/permission"
)

const authCookieName = "token"

type HashService interface {
	Hash(string) (Hash, Salt, error)
	HashAndCompare(string, Hash, Salt) bool
}

type TokenService interface {
	SignWithClaims(Claims) (string, error)
	Validate(string) (*Claims, error)
}

type Claims struct {
	jwt.RegisteredClaims
	RoleName    string                  `json:"role"`
	Permissions []permission.Permission `json:"permissions"`
}

type AuthStore interface {
	RoleByUsername(string) (role.Role, error)
	PasswordHashAndSaltByUsername(string) (Hash, Salt, error)
}

type (
	Hash []byte
	Salt []byte
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func LoginHandler(
	authStore AuthStore,
	tokenService TokenService,
	hashService HashService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hash, salt, err := authStore.PasswordHashAndSaltByUsername(creds.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !hashService.HashAndCompare(creds.Password, hash, salt) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userRole, err := authStore.RoleByUsername(creds.Username)
		if err != nil {
			http.Error(w, "you don't have a role", http.StatusConflict)
			return
		}

		now := time.Now()
		expiresAt := now.Add(24 * time.Hour)
		claims := Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "mdma-backend",
				Subject:   creds.Username,
			},
			RoleName:    userRole.Name,
			Permissions: userRole.Permissions,
		}

		tokenStr, err := tokenService.SignWithClaims(claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:     authCookieName,
			Value:    tokenStr,
			Expires:  expiresAt,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		token := Token{
			Token:     tokenStr,
			ExpiresAt: expiresAt,
		}
		render.JSON(w, r, token)
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &http.Cookie{
			Name:     authCookieName,
			Value:    "",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}
}
