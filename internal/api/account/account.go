package account

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

	r.Route("/accounts", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", s.getAccountUser())
			r.Get("/", nil)
			r.Post("/", nil)
			r.Put("/{id}", nil)
			r.Delete("/{id}", nil)
		})
		r.Route("/services", func(r chi.Router) {
			r.Get("/{id}", nil)
			r.Get("/", nil)
			r.Post("/", nil)
			r.Put("/{id}", nil)
			r.Delete("/{id}", nil)
		})
	})

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAccountUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}
