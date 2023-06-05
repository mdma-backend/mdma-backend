package mesh_node

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mdma-backend/mdma-backend/internal/api/data"
	"net/http"
	"strconv"
)

type MeshNode struct {
	Uuid      string
	Latitude  float32
	Longitude float32
	CreatedAt string
	UpdatedAt string
	UpdateId  float32
}

type MeshNodeStore interface {
	GetMeshNodes() ([]MeshNode, error)
	GetMeshNode(uuid string) (MeshNode, error)
	PostMeshNode(latitude float32, longitude float32, updateId float32) error
	PostMeshNodeData(controllerUuid string, meshNodeType string, value string, measuredAt string) error
	PutMeshNode(uuid string, latitude float32, longitude float32) error
	DeleteMeshNode(uuid string) error
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
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method not allowed"))

			return
		}

		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) getMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method not allowed"))

			return
		}

		// Extrahiere die UUID aus dem Request-URL-Pfad
		uuid := chi.URLParam(r, "uuid")

		// Daten aus der Datenbank abrufen
		meshNode, err := s.meshNodeStore.GetMeshNode(uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// JSON-Antwort erstellen
		response, err := json.Marshal(meshNode)
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

func (s service) postMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method not allowed"))

			return
		}

		// Dekodiere den JSON-Body der Anfrage in ein Payload-Objekt
		var meshNode MeshNode
		err := json.NewDecoder(r.Body).Decode(&meshNode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Invalid request payload"))

			return
		}

		err = s.meshNodeStore.PostMeshNode(meshNode.Latitude, meshNode.Longitude, meshNode.UpdateId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// Sende eine Erfolgsantwort zurück
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST request successful"))
	}
}

func (s service) postMeshNodeData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method not allowed"))

			return
		}

		// Dekodiere den JSON-Body der Anfrage in ein Payload-Objekt
		var meshNodeData data.Data
		err := json.NewDecoder(r.Body).Decode(&meshNodeData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Invalid request payload"))

			return
		}

		uuid := chi.URLParam(r, "uuid")
		err = s.meshNodeStore.PostMeshNodeData(uuid, meshNodeData.Type, meshNodeData.Value, meshNodeData.MeasuredAt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// Sende eine Erfolgsantwort zurück
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST request successful"))
	}
}

func (s service) putMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method not allowed"))

			return
		}

		// Extrahiere die UUID aus dem Request-URL-Pfad
		uuid := chi.URLParam(r, "uuid")
		latString := r.URL.Query().Get("lat")
		lngString := r.URL.Query().Get("lng")

		lat, _ := strconv.ParseFloat(latString, 32)
		lng, _ := strconv.ParseFloat(lngString, 32)

		err := s.meshNodeStore.PutMeshNode(uuid, float32(lat), float32(lng))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("501 not implemented"))
	}
}

func (s service) deleteMeshNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Method not allowed"))

			return
		}

		uuid := chi.URLParam(r, "uuid")

		// Daten löschen
		err := s.meshNodeStore.DeleteMeshNode(uuid)
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
