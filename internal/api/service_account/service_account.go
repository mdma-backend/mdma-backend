package service_account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type ServiceUserStore interface {
	ServiceAccount(id int) (ServiceAccount, error)
	AllServiceAccounts() ([]ServiceAccount, error)
	CreateServiceAccount(roleID int, username string) error
	UpdateServiceAccount(id int, roleID int, username string) (ServiceAccount, error)
	DeleteServiceAccount(id int) error
}

type service struct {
	handler          http.Handler
	serviceUserStore ServiceUserStore
}
type ServiceAccount struct {
	ID        int    `json:"id,omitempty"`
	RoleID    int    `json:"roleId,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Username  string `json:"username,omitempty"`
	Token     []byte `json:"token,omitempty"`
}

func NewService(serviceUserStore ServiceUserStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:          r,
		serviceUserStore: serviceUserStore,
	}

	r.Get("/{id}", auth.RestrictHandlerFunc(s.getAccountService(), permission.ServiceAccountRead))
	r.Get("/", auth.RestrictHandlerFunc(s.getAllService(), permission.ServiceAccountRead))
	r.Post("/", auth.RestrictHandlerFunc(s.createAccountService(), permission.ServiceAccountCreate))
	r.Put("/{id}", auth.RestrictHandlerFunc(s.updateAccountService(), permission.ServiceAccountUpdate))
	r.Delete("/{id}", auth.RestrictHandlerFunc(s.deleteAccountService(), permission.ServiceAccountDelete))

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAllService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.serviceUserStore.AllServiceAccounts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
			return
		}

		response, err := json.Marshal(accounts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (s service) getAccountService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)

		account, err := s.serviceUserStore.ServiceAccount(idInt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid account ID")
			return
		}

		response, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (s service) createAccountService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serviceAccount ServiceAccount

		if err := json.NewDecoder(r.Body).Decode(&serviceAccount); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid request payload")
			return
		}

		err := s.serviceUserStore.CreateServiceAccount(serviceAccount.RoleID, serviceAccount.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to create User")
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serviceAccount)
	}
}
func (s service) updateAccountService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, _ := strconv.Atoi(id)

		var serviceAccount ServiceAccount
		err := json.NewDecoder(r.Body).Decode(&serviceAccount)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid request payload")
			return
		}

		serviceAccount, err = s.serviceUserStore.UpdateServiceAccount(idInt, serviceAccount.RoleID, serviceAccount.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to update user")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "UserAccount updated successfully")
	}
}
func (s service) deleteAccountService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.serviceUserStore.DeleteServiceAccount(idInt); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to delete service Account")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "UserAccount deleted successfully")
	}
}
