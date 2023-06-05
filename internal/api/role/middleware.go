package role

import (
	"context"
	"net/http"

	"github.com/mdma-backend/mdma-backend/internal/api/role/permission"
)

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
	return nil
}
