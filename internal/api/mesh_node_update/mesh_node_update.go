package mesh_node_update

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type MeshNodeUpdateStore interface {
	MeshNodeUpdateByID(types.MeshNodeUpdateID) (types.MeshNodeUpdate, error)
	MeshNodeUpdates() ([]types.MeshNodeUpdate, error)
	CreateMeshNodeUpdate(*types.MeshNodeUpdate) error
}

type service struct {
	handler             http.Handler
	meshNodeUpdateStore MeshNodeUpdateStore
}

func NewService(meshNodeUpdateStore MeshNodeUpdateStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:             r,
		meshNodeUpdateStore: meshNodeUpdateStore,
	}

	r.Get("/", auth.RestrictHandlerFunc(s.getMeshNodeUpdates(), permission.MeshNodeUpdateRead))
	r.Get("/{id}", auth.RestrictHandlerFunc(s.getMeshNodeUpdate(), permission.MeshNodeUpdateRead))
	r.Post("/", auth.RestrictHandlerFunc(s.postMeshNodeUpdate(), permission.MeshNodeUpdateCreate))

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getMeshNodeUpdates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		meshNodeUpdates, err := s.meshNodeUpdateStore.MeshNodeUpdates()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, meshNodeUpdates)
	}
}

func (s service) getMeshNodeUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		meshNodeUpdateID, err := types.IDFromString[types.MeshNodeUpdateID](id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		meshNodeUpdate, err := s.meshNodeUpdateStore.MeshNodeUpdateByID(meshNodeUpdateID)
		if errors.Is(err, types.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, meshNodeUpdate)
	}
}

func (s service) postMeshNodeUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var meshNodeUpdate types.MeshNodeUpdate
		if err := json.NewDecoder(r.Body).Decode(&meshNodeUpdate); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.meshNodeUpdateStore.CreateMeshNodeUpdate(&meshNodeUpdate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, meshNodeUpdate)
	}
}
