package user_account

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type UserStore interface {
	UserAccount(types.UserAccountID) (types.UserAccount, error)
	AllUserAccounts() ([]types.UserAccount, error)
	CreateUserAccount(*types.UserAccount, types.Hash, types.Salt) error
	UpdateUserAccount(types.UserAccountID, *types.UserAccount) error
	DeleteUserAccount(types.UserAccountID) error
}

type service struct {
	handler     http.Handler
	userStore   UserStore
	hashService types.HashService
}

func NewService(userStore UserStore, hashService types.HashService) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:     r,
		userStore:   userStore,
		hashService: hashService,
	}

	r.Get("/{id}", auth.RestrictHandlerFunc(s.getAccountUser(), permission.UserAccountRead))
	r.Get("/", auth.RestrictHandlerFunc(s.getAllUsers(), permission.UserAccountRead))
	r.Post("/", auth.RestrictHandlerFunc(s.createAccountUser(), permission.UserAccountCreate))
	r.Put("/{id}", auth.RestrictHandlerFunc(s.updateAccountUser(), permission.UserAccountUpdate))
	r.Delete("/{id}", auth.RestrictHandlerFunc(s.deleteAccountUser(), permission.UserAccountDelete))

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAccounts, err := s.userStore.AllUserAccounts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, userAccounts)
	}
}

func (s service) getAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		userAccountID, err := types.IDFromString[types.UserAccountID](id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userAccount, err := s.userStore.UserAccount(userAccountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, userAccount)
	}
}

func (s service) createAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userAccount types.UserAccount
		if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hash, salt, err := s.hashService.Hash(userAccount.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := s.userStore.CreateUserAccount(&userAccount, hash, salt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Clear the password field for safety?
		userAccount.Password = ""

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, userAccount)
	}
}

func (s service) updateAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		userAccountID, err := types.IDFromString[types.UserAccountID](id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var userAccount types.UserAccount
		if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = s.userStore.UpdateUserAccount(userAccountID, &userAccount); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, userAccount)
	}
}

func (s service) deleteAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		userAccountID, err := types.IDFromString[types.UserAccountID](id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.userStore.DeleteUserAccount(userAccountID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
