package me

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types"
)

type UserAccountStore interface {
	UserAccountByID(types.UserAccountID) (types.UserAccount, error)
}

type ServiceAccountStore interface {
	ServiceAccountByID(types.ServiceAccountID) (types.ServiceAccount, error)
}

type service struct {
	handler             http.Handler
	userAccountStore    UserAccountStore
	serviceAccountStore ServiceAccountStore
}

func NewService(userAccountStore UserAccountStore, serviceAccountStore ServiceAccountStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:             r,
		userAccountStore:    userAccountStore,
		serviceAccountStore: serviceAccountStore,
	}

	r.Get("/", s.getMe())

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, ok := r.Context().Value(auth.AccountInfoCtxKey).(types.AccountInfo)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		render.JSON(w, r, info)
	}
}
