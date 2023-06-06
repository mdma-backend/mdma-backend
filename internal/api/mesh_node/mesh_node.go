package mesh_node

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mdma-backend/mdma-backend/internal/api/data"
	"net/http"
)

type Point struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type MeshNode struct {
	Id        string `json:"uuid"`
	UpdateId  int64  `json:"updateId,omitempty"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Location  Point  `json:"location"`
}

type MeshNodeStore interface {
	MeshNodes() ([]MeshNode, error)
	MeshNodeById(id string) (MeshNode, error)
	CreateMeshNode(node MeshNode) error
	CreateMeshNodeData(meshNodeId string, data data.Data) error
	UpdateMeshNode(id string, node MeshNode) error
	DeleteMeshNode(id string) error
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

		meshNodes, err := s.meshNodeStore.MeshNodes()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// JSON-Antwort erstellen
		response, err := json.Marshal(meshNodes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 JSON conversion failed"))
			return
		}

		// Antwort senden
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
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
		meshNode, err := s.meshNodeStore.MeshNodeById(uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		if meshNode.Id == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 mesh node not found"))
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

		err = s.meshNodeStore.CreateMeshNode(meshNode)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// Sende eine Erfolgsantwort zurück
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mesh node - POST request successful"))
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
		err = s.meshNodeStore.CreateMeshNodeData(uuid, meshNodeData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// Sende eine Erfolgsantwort zurück
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mesh node data - POST request successful"))
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
		var meshNode MeshNode

		if err := json.NewDecoder(r.Body).Decode(&meshNode); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}
		/*
			updateId := r.Body.("updateId")
			latString := r.URL.Query().Get("latitude")
			lonString := r.URL.Query().Get("longitude")

			lat, _ := strconv.ParseFloat(latString, 32)
			lon, _ := strconv.ParseFloat(lonString, 32)
			updateIdInt, _ := strconv.ParseInt(updateId, 10, 64)
			location := Point{Lat: float32(lat), Lon: float32(lon)}
		*/
		err := s.meshNodeStore.UpdateMeshNode(uuid, meshNode)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		// Erfolgsmeldung senden
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mesh node - PUT request successful"))
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
		w.Write([]byte("Mesh node - DELETE request successful"))
	}
}
