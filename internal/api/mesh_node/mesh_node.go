package mesh_node

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/api/data"
	"github.com/mdma-backend/mdma-backend/internal/types"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type MeshNodeStore interface {
	MeshNodes() ([]types.MeshNode, error)
	MeshNodeById(types.UUID) (types.MeshNode, error)
	CreateMeshNode(*types.MeshNode) error
	CreateMeshNodeData(types.UUID, *data.Data) error
	CreateManyMeshNodeData(types.UUID, []data.Data) error
	UpdateMeshNode(types.UUID, *types.MeshNode) error
	DeleteMeshNode(types.UUID) error
}

type service struct {
	handler       http.Handler
	meshNodeStore MeshNodeStore
}

func NewService(store MeshNodeStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:       r,
		meshNodeStore: store,
	}

	r.Get("/", auth.RestrictHandlerFunc(s.getMeshNodes(), permission.MeshNodeRead))
	r.Get("/{uuid}", auth.RestrictHandlerFunc(s.getMeshNode(), permission.MeshNodeRead))
	r.Post("/", auth.RestrictHandlerFunc(s.postMeshNode(), permission.MeshNodeCreate))
	r.Post("/{uuid}/data", auth.RestrictHandlerFunc(s.postMeshNodeData(), permission.DataCreate))
	r.Post("/{uuid}/data-list", auth.RestrictHandlerFunc(s.postManyMeshNodeData(), permission.DataCreate))
	r.Put("/{uuid}", auth.RestrictHandlerFunc(s.putMeshNode(), permission.MeshNodeUpdate))
	r.Delete("/{uuid}", auth.RestrictHandlerFunc(s.deleteMeshNode(), permission.MeshNodeDelete))

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getMeshNodes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		meshNodes, err := s.meshNodeStore.MeshNodes()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, meshNodes)
	}
}

func (s service) getMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := chi.URLParam(r, "uuid")
		meshNodeUUID, err := types.UUIDFromString(uuidStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		meshNode, err := s.meshNodeStore.MeshNodeById(meshNodeUUID)
		if errors.Is(err, types.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, meshNode)
	}
}

func (s service) postMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var meshNode types.MeshNode
		if err := json.NewDecoder(r.Body).Decode(&meshNode); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.meshNodeStore.CreateMeshNode(&meshNode); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, meshNode)
	}
}

func (s service) postMeshNodeData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := chi.URLParam(r, "uuid")
		meshNodeUUID, err := types.UUIDFromString(uuidStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var data data.Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data.MeshNodeUUID = meshNodeUUID.String()

		if err = s.meshNodeStore.CreateMeshNodeData(meshNodeUUID, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, data)
	}
}

func (s service) postManyMeshNodeData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := chi.URLParam(r, "uuid")
		meshNodeUUID, err := types.UUIDFromString(uuidStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var meshNodeData []data.Data
		if err := json.NewDecoder(r.Body).Decode(&meshNodeData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = s.meshNodeStore.CreateManyMeshNodeData(meshNodeUUID, meshNodeData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, meshNodeData)
	}
}

func (s service) putMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := chi.URLParam(r, "uuid")
		meshNodeUUID, err := types.UUIDFromString(uuidStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var meshNode types.MeshNode
		if err := json.NewDecoder(r.Body).Decode(&meshNode); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.meshNodeStore.UpdateMeshNode(meshNodeUUID, &meshNode); errors.Is(err, types.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		meshNode.UUID = meshNodeUUID

		render.JSON(w, r, meshNode)
	}
}

func (s service) deleteMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := chi.URLParam(r, "uuid")
		meshNodeUUID, err := types.UUIDFromString(uuidStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.meshNodeStore.DeleteMeshNode(meshNodeUUID); errors.Is(err, types.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
