package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mdma-backend/mdma-backend/internal/types"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

var (
	AccountInfoCtxKey = &struct{}{}
)

type RoleStore interface {
	RoleByUserAccountID(uaId types.UserAccountID) (types.Role, error)
	RoleByServiceAccountID(saId types.ServiceAccountID) (types.Role, error)
}

func Middleware(
	tokenService types.TokenService,
	roleStore RoleStore,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := claimsFromRequest(r, tokenService)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			var role types.Role
			switch claims.AccountType {
			case types.UserAccountType:
				role, err = roleStore.RoleByUserAccountID(types.UserAccountID(claims.AccountID))
			case types.ServiceAccountType:
				role, err = roleStore.RoleByServiceAccountID(types.ServiceAccountID(claims.AccountID))
			default:
				http.Error(w, "invalid account type", http.StatusBadRequest)
				return
			}
			if err != nil {
				http.Error(w, "account has no role", http.StatusBadRequest)
				return
			}

			info := types.AccountInfo{
				AccountType: claims.AccountType,
				AccountID:   claims.AccountID,
				Role:        role,
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, AccountInfoCtxKey, info)
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
				http.Error(w, fmt.Sprintf("missing permission %s", permission), http.StatusUnauthorized)
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

func claimsFromRequest(r *http.Request, tokenService types.TokenService) (*types.Claims, error) {
	var tokenStr string
	if bearerStr := r.Header.Get("Authorization"); bearerStr != "" {
		tokenStr = strings.TrimPrefix(bearerStr, "Bearer ")
	} else {
		cookie, err := r.Cookie(authCookieName)
		if err != nil {
			return nil, errors.New("no token in header or cookie")
		}

		tokenStr = cookie.Value
	}

	claims, err := tokenService.Validate(tokenStr)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func permissionsFromContext(ctx context.Context) []permission.Permission {
	info, ok := ctx.Value(AccountInfoCtxKey).(types.AccountInfo)
	if !ok {
		return nil
	}

	return info.Role.Permissions
}
