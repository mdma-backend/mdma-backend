package area

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type service struct {
	handler http.Handler
}

type Area struct {
	AreaId      int
	MeshNodeIds []string
}

var areas = map[int]Area{
	1: {AreaId: 1, MeshNodeIds: []string{"a53b3f71-f073-4578-9557-92fd19d93bb9", "c33ea7b6-68a7-4bc6-b1e9-0c365db74081", "f1aef837-04ac-4316-ae1f-0465bc2eb2fa"}},
	2: {AreaId: 2, MeshNodeIds: []string{"f1aef837-04ac-4316-ae1f-0465bc2eb2fa", "a8957622-acc5-4ddb-bb1f-17e63d3a514f"}},
	3: {AreaId: 3, MeshNodeIds: []string{"a53b3f71-f073-4578-9557-92fd19d93bb9", "c33ea7b6-68a7-4bc6-b1e9-0c365db74081", "f1aef837-04ac-4316-ae1f-0465bc2eb2fa", "a8957622-acc5-4ddb-bb1f-17e63d3a514f"}},
}

func NewService() http.Handler {
	r := chi.NewRouter()
	s := service{
		handler: r,
	}

	r.Get("/", auth.RestrictHandlerFunc(s.getAreas(), permission.AreaRead))
	r.Get("/{id}", auth.RestrictHandlerFunc(s.getArea(), permission.AreaRead))

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getArea() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idParam)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("204 invalid ID"))
			return
		}

		area, ok := areas[id]
		if !ok {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("204 internal server error"))
			return
		}

		response, err := json.Marshal(area)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (s service) getAreas() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := json.Marshal(areas)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
