package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserStore interface {
	User(id int) (UserAccount, error)
	PostUser(roleID int, username string, password []byte) (int, error)
	PutUser(id int, username string, password []byte) (UserAccount, error)
	DeleteUser(id int) error
}

type service struct {
	handler   http.Handler
	userStore UserStore
}
type UserAccount struct {
	ID        int
	RoleID    int
	CreatedAt string
	UpdatedAt string
	Username  string
	Password  []byte
}

func NewService(userStore UserStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:   r,
		userStore: userStore,
	}

	r.Get("/{id}", s.getAccountUser())
	r.Post("/", s.createAccountUser())
	r.Put("/{id}", s.updateAccountUser())
	r.Delete("/{id}", s.deleteAccountUser())

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, _ := strconv.Atoi(id) // Convert Account-ID-String in Integer

		//Daten holen
		account, err := s.userStore.User(idInt)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Account not found")
			return
		}

		// JSON-Antwort erstellen
		response, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
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
		var account UserAccount
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid request payload")
			return
		}

		// JSON-Antwort senden
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(account)
	}
}
func (s service) updateAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, _ := strconv.Atoi(id)

		var account UserAccount
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid request payload")
			return
		}

		// Datenbank aufrufen, um den Benutzer zu aktualisieren
		account, err = s.userStore.PutUser(idInt, account.Username, account.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to update user")
			return
		}

		// Erfolgreiche Antwort senden
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User updated successfully")
	}

}
func (s service) deleteAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, _ := strconv.Atoi(id)

		// Delete the user account
		err := s.userStore.DeleteUser(idInt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to delete user")
			return
		}

		// Successful response
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User deleted successfully")
	}
}
