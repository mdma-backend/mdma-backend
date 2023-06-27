package area

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mdma-backend/mdma-backend/internal/types"
)

type Area struct {
	AreaID        int      `json:"areaId"`
	MeshNodeUUIDs []string `json:"meshNodeUUIDs"`
}

var areas = map[uint]Area{
	1: {
		AreaID: 1,
		MeshNodeUUIDs: []string{
			"a53b3f71-f073-4578-9557-92fd19d93bb9",
			"c33ea7b6-68a7-4bc6-b1e9-0c365db74081",
			"f1aef837-04ac-4316-ae1f-0465bc2eb2fa",
		},
	},
	2: {
		AreaID: 2,
		MeshNodeUUIDs: []string{
			"f1aef837-04ac-4316-ae1f-0465bc2eb2fa",
			"a8957622-acc5-4ddb-bb1f-17e63d3a514f",
		},
	},
	3: {
		AreaID: 3,
		MeshNodeUUIDs: []string{
			"a53b3f71-f073-4578-9557-92fd19d93bb9",
			"c33ea7b6-68a7-4bc6-b1e9-0c365db74081",
			"f1aef837-04ac-4316-ae1f-0465bc2eb2fa",
			"a8957622-acc5-4ddb-bb1f-17e63d3a514f",
		},
	},
}

type service struct {
	handler http.Handler
}

func NewService() http.Handler {
	r := chi.NewRouter()
	s := service{
		handler: r,
	}

	r.Get("/", s.getAreas())
	r.Get("/{id}", s.getArea())

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getArea() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := types.IDFromString[uint](idParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		area, ok := areas[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		render.JSON(w, r, area)
	}
}

func (s service) getAreas() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var aa []Area
		for _, a := range areas {
			aa = append(aa, a)
		}

		render.JSON(w, r, aa)
	}
}
