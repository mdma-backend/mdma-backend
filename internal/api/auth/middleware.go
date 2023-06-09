package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

var (
	PermissionsCtxKey = &struct{}{}
)

func Middleware(tokenService TokenService, excludedPaths ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range excludedPaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			var tokenStr string
			if bearerStr := r.Header.Get("Authorization"); bearerStr != "" {
				tokenStr = strings.TrimPrefix(bearerStr, "Bearer ")
			} else {
				cookie, err := r.Cookie(authCookieName)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				tokenStr = cookie.Value
			}

			claims, err := tokenService.Validate(tokenStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, PermissionsCtxKey, claims.Permissions)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func RestrictHandlerFunc(next http.HandlerFunc, permissions ...permission.Permission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		perms := permissionsFromContext(r.Context())

		for _, permission := range permissions {
			var hasPermission bool
			for _, p := range perms {
				if p == permission {
					hasPermission = true
					continue
				}
			}

			if !hasPermission {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("missing permission " + permission))
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}

func RestrictHandler(next http.Handler, permissions ...permission.Permission) http.Handler {
	return http.HandlerFunc(RestrictHandlerFunc(next.ServeHTTP, permissions...))
}

func RestrictMiddleware(permissions ...permission.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RestrictHandler(next)
	}
}

func permissionsFromContext(ctx context.Context) []permission.Permission {
	return ctx.Value(PermissionsCtxKey).([]permission.Permission)
}
