package data

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DataStore interface {
	DeleteData(uuid string) error
	Types() ([]string, error)
}

type service struct {
	handler   http.Handler
	dataStore DataStore
}

func NewService(dataStore DataStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:   r,
		dataStore: dataStore,
	}

	r.Get("/", s.getManyData())
	r.Get("/{uuid}", s.getData())
	r.Get("/types", s.getDataTypes())
	r.Get("/aggregated", s.getAggregatedData())
	r.Delete("/{uuid}", s.deleteData())

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAggregatedData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) getManyData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) getData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//del
		//wir nutzen chi	https://github.com/go-chi/chi
		//für query parameter --> chi.URLParam(r, "uuid") gibt UUID wieder aus der Request
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) getDataTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataTypes, err := s.dataStore.Types()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// JSON-Antwort erstellen
		response, err := json.Marshal(dataTypes)
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

func (s service) deleteData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := chi.URLParam(r, "uuid")

		// Daten löschen
		err := s.dataStore.DeleteData(uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// Erfolgsmeldung senden
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data deleted successfully"))
	}
}
