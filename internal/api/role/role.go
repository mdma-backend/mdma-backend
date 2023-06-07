package role

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type ErrNotFound struct {
	Err    error
	RoleID types.RoleID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("role with id %d not found: %s", e.RoleID, e.Err)
}

type RoleStore interface {
	RoleByID(types.RoleID) (types.Role, error)
	Roles() ([]types.Role, error)
	CreateRole(*types.Role) error
	UpdateRole(types.RoleID, *types.Role) error
	DeleteRole(types.RoleID) error
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

	r.Get("/", auth.RestrictHandlerFunc(s.getRoles(), permission.RoleRead))
	r.Get("/{id}", auth.RestrictHandlerFunc(s.getRole(), permission.RoleRead))
	r.Post("/", auth.RestrictHandlerFunc(s.postRole(), permission.RoleCreate))
	r.Put("/{id}", auth.RestrictHandlerFunc(s.putRole(), permission.RoleUpdate))
	r.Delete("/{id}", auth.RestrictHandlerFunc(s.deleteRole(), permission.RoleDelete))

	r.Get("/permissions", auth.RestrictHandlerFunc(getPermissions(), permission.RoleCreate))

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
		roleID, err := types.RoleIDFromString(id)
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
		var role types.Role
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
		roleID, err := types.RoleIDFromString(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var role types.Role
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
		roleID, err := types.RoleIDFromString(id)
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

func getPermissions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, permission.Permissions())
	}
}
