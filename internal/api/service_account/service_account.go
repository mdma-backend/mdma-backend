package service_account

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type ServiceAccountStore interface {
	RoleByID(types.RoleID) (types.Role, error)
	ServiceAccount(types.ServiceAccountID) (types.ServiceAccount, error)
	AllServiceAccounts() ([]types.ServiceAccount, error)
	CreateServiceAccount(*types.ServiceAccount) error
	UpdateServiceAccount(types.ServiceAccountID, *types.ServiceAccount) error
	UpdateServiceAccountToken(types.ServiceAccountID, types.Token) error
	DeleteServiceAccount(types.ServiceAccountID) error
}

type service struct {
	handler             http.Handler
	serviceAccountStore ServiceAccountStore
	tokenService        types.TokenService
}

func NewService(serviceUserStore ServiceAccountStore, tokenSerive types.TokenService) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:             r,
		serviceAccountStore: serviceUserStore,
		tokenService:        tokenSerive,
	}

	r.Get("/{id}", auth.RestrictHandlerFunc(s.getAccountService(), permission.ServiceAccountRead))
	r.Get("/", auth.RestrictHandlerFunc(s.getAllService(), permission.ServiceAccountRead))
	r.Post("/", auth.RestrictHandlerFunc(s.createAccountService(), permission.ServiceAccountCreate))
	r.Post("/{id}/refresh-token", auth.RestrictHandlerFunc(s.refreshAccountServiceToken(), permission.ServiceAccountCreate))
	r.Put("/{id}", auth.RestrictHandlerFunc(s.updateAccountService(), permission.ServiceAccountUpdate))
	r.Delete("/{id}", auth.RestrictHandlerFunc(s.deleteAccountService(), permission.ServiceAccountDelete))

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAllService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAccounts, err := s.serviceAccountStore.AllServiceAccounts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, userAccounts)
	}
}

func (s service) getAccountService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		serviceAccountID, err := types.IDFromString[types.ServiceAccountID](id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		serviceAccount, err := s.serviceAccountStore.ServiceAccount(serviceAccountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, serviceAccount)
	}
}

func (s service) createAccountService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serviceAccount types.ServiceAccount
		if err := json.NewDecoder(r.Body).Decode(&serviceAccount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.serviceAccountStore.CreateServiceAccount(&serviceAccount); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		now := time.Now()
		expiresAt := now.Add(24 * 7 * 52 * time.Hour) // one year
		claims := types.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "mdma-backend",
				Subject:   serviceAccount.Name,
			},
			AccountType: types.ServiceAccountType,
			AccountID:   uint(serviceAccount.ID),
		}

		token, err := s.tokenService.SignWithClaims(claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		serviceAccount.Token = token.Value

		if err := s.serviceAccountStore.UpdateServiceAccountToken(serviceAccount.ID, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, serviceAccount)
	}
}

func (s service) refreshAccountServiceToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		serviceAccountID, err := types.IDFromString[types.ServiceAccountID](id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		serviceAccount, err := s.serviceAccountStore.ServiceAccount(serviceAccountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		now := time.Now()
		expiresAt := now.Add(24 * 7 * 52 * time.Hour) // one year
		claims := types.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiresAt),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "mdma-backend",
				Subject:   serviceAccount.Name,
			},
			AccountType: types.ServiceAccountType,
			AccountID:   uint(serviceAccount.ID),
		}

		token, err := s.tokenService.SignWithClaims(claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		serviceAccount.Token = token.Value

		if err = s.serviceAccountStore.UpdateServiceAccountToken(serviceAccountID, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, serviceAccount)
	}
}

func (s service) updateAccountService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		serviceAccountID, err := types.IDFromString[types.ServiceAccountID](id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var serviceAccount types.ServiceAccount
		if err := json.NewDecoder(r.Body).Decode(&serviceAccount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = s.serviceAccountStore.UpdateServiceAccount(serviceAccountID, &serviceAccount); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, serviceAccount)
	}
}

func (s service) deleteAccountService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		serviceAccountID, err := types.IDFromString[types.ServiceAccountID](id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.serviceAccountStore.DeleteServiceAccount(serviceAccountID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
