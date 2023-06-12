package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mdma-backend/mdma-backend/internal/types"
)

const authCookieName = "token"

type AuthStore interface {
	RoleByUsername(string) (types.Role, error)
	PasswordHashAndSaltByUsername(string) (types.Hash, types.Salt, error)
}

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
	tokenService types.TokenService,
	hashService types.HashService,
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
		claims := types.Claims{
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

		token, err := tokenService.SignWithClaims(claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:     authCookieName,
			Value:    token.Value,
			Expires:  expiresAt,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

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
