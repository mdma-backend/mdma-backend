package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserStore interface {
	UserAccount(id int) (UserAccount, error)
	AllUserAccounts() ([]UserAccount, error)
	CreateUserAccount(roleID int, createdAt string, username string, password []byte) error
	UpdateUserAccount(id int, roleID int, username string, password []byte) (UserAccount, error)
	DeleteUserAccount(id int) error
}

type service struct {
	handler   http.Handler
	userStore UserStore
}
type UserAccount struct {
	ID        int    `json:"id,omitempty"`
	RoleID    int    `json:"roleId,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  []byte `json:"password,omitempty"`
}

func NewService(userStore UserStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:   r,
		userStore: userStore,
	}

	r.Route("/users", func(r chi.Router) {
		r.Get("/{id}", s.getAccountUser())
		r.Get("/", s.getAllUsers())
		r.Post("/", s.createAccountUser())
		r.Put("/{id}", s.updateAccountUser())
		r.Delete("/{id}", s.deleteAccountUser())
	})

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Alle Benutzerkonten abrufen
		accounts, err := s.userStore.AllUserAccounts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
			return
		}

		// JSON-Antwort erstellen
		response, err := json.Marshal(accounts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
			return
		}

		// Antwort senden
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (s service) getAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id) // Convert Account-ID-String in Integer

		//Daten holen
		account, err := s.userStore.UserAccount(idInt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid account ID")
			return
		}

		// JSON-Antwort erstellen
		response, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
			return
		}

		// Antwort senden
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
func (s service) createAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userAccount UserAccount

		if err := json.NewDecoder(r.Body).Decode(&userAccount); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid request payload")
			return
		}

		err := s.userStore.CreateUserAccount(userAccount.RoleID, userAccount.CreatedAt, userAccount.Username, userAccount.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to create User")
			return
		}

		// JSON-Antwort senden
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userAccount)
	}
}
func (s service) updateAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, _ := strconv.Atoi(id)

		var userAccount UserAccount
		err := json.NewDecoder(r.Body).Decode(&userAccount)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid request payload")
			return
		}

		// Datenbank aufrufen, um den Benutzer zu aktualisieren
		userAccount, err = s.userStore.UpdateUserAccount(idInt, userAccount.RoleID, userAccount.Username, userAccount.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to update user")
			return
		}

		// Erfolgreiche Antwort senden
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "UserAccount updated successfully")
	}
}
func (s service) deleteAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Delete the user account
		if err := s.userStore.DeleteUserAccount(idInt); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to delete user")
			return
		}

		// Successful response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "UserAccount deleted successfully")
	}
}
