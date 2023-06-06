package role

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mdma-backend/mdma-backend/internal/api/role/permission"
)

type ErrNotFound struct {
	Err    error
	RoleID RoleID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("role with id %d not found: %s", e.RoleID, e.Err)
}

type RoleStore interface {
	RoleByID(RoleID) (Role, error)
	Roles() ([]Role, error)
	CreateRole(*Role) error
	UpdateRole(RoleID, *Role) error
	DeleteRole(RoleID) error
}

type RoleID int

func IDFromString(s string) (RoleID, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%s is not an integer", s)
	}

	if id <= 0 {
		return 0, errors.New("id must be greater than 0")
	}

	return RoleID(id), nil
}

type Role struct {
	ID          RoleID                  `json:"id,omitempty"`
	CreatedAt   time.Time               `json:"createAt"`
	UpdatedAt   *time.Time              `json:"updatedAt,omitempty"`
	Name        string                  `json:"name"`
	Permissions []permission.Permission `json:"permissions,omitempty"`
}

type service struct {
	handler   http.Handler
	roleStore RoleStore
}

func NewService(roleStore RoleStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:   r,
		roleStore: roleStore,
	}

	r.Get("/", s.getRoles())
	r.Get("/{id}", s.getRole())
	r.Post("/", s.postRole())
	r.Put("/{id}", s.putRole())
	r.Delete("/{id}", s.deleteRole())

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getRoles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles, err := s.roleStore.Roles()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, roles)
	}
}

func (s service) getRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		roleID, err := IDFromString(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		role, err := s.roleStore.RoleByID(roleID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, role)
	}
}

func (s service) postRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var role Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.roleStore.CreateRole(&role); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, role)
	}
}

func (s service) putRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		roleID, err := IDFromString(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var role Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := s.roleStore.UpdateRole(roleID, &role); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		role.ID = roleID

		render.JSON(w, r, role)
	}
}

func (s service) deleteRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		roleID, err := IDFromString(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.roleStore.DeleteRole(roleID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
