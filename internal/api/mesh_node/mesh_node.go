package mesh_node

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type service struct {
	handler http.Handler
}

func NewService() http.Handler {
	r := chi.NewRouter()
	s := service{
		handler: r,
	}

	r.Get("/", s.getMeshNodes())
	r.Get("/{uuid}", s.getMeshNode())
	r.Post("/", s.postMeshNode())
	r.Post("/{uuid}/data", s.postMeshNodeData())
	r.Put("/{uuid}", s.putMeshNode())
	r.Delete("/{uuid}", s.deleteMeshNode())

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getMeshNodes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) getMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) postMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) postMeshNodeData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) putMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) deleteMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}
