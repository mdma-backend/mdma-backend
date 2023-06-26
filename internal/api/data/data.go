package data

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type DataStore interface {
	GetAggregatedData(dataType string, meshNodeUUIDs []string, startTime time.Time, endTime time.Time, sampleTime time.Duration, sampleCount int, aggregateFunction string) (AggregatedData, error)
	GetManyData(dataType string, meshNodeUUIDs []string, startTime time.Time, endTime time.Time) (ManyData, error)
	GetData(uuid string) (Data, error)
	DeleteData(uuid string) error
	GetTypes() ([]string, error)
}

type service struct {
	handler   http.Handler
	dataStore DataStore
}

type AggregatedData struct {
	AggregateFunction string   `json:"aggregationFunction"`
	DataType          string   `json:"type"`
	MeshNodeUUIDs     []string `json:"meshNodeUUIDs"`
	Samples           []Sample `json:"samples"`
}

type Sample struct {
	IntervalStartAt string `json:"intervalStartAt"`
	IntervalEndAt   string `json:"intervalEndAt"`
	Value           string `json:"value"`
}

type ManyData struct {
	DataType      string         `json:"type"`
	MeasuredDatas []MeasuredData `json:"data"`
}

type MeasuredData struct {
	MeshnodeUUID string        `json:"meshNodeUUID"`
	Measurements []Measurement `json:"measurements"`
}

type Measurement struct {
	UUID       string `json:"UUID"`
	MeasuredAt string `json:"measuredAt"`
	Value      string `json:"value"`
}

type Data struct {
	UUID         string `json:"uuid"`
	MeshNodeUUID string `json:"meshNodeUUID"`
	Type         string `json:"type"`
	CreatedAt    string `json:"createdAt"`
	MeasuredAt   string `json:"measuredAt"`
	Value        string `json:"value"`
}

func NewService(dataStore DataStore, tokenService types.TokenService, roleService auth.RoleStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:   r,
		dataStore: dataStore,
	}

	r.Get("/", s.getManyData())
	r.Get("/{uuid}", s.getData())
	r.Get("/types", s.getDataTypes())
	r.Get("/aggregated", s.getAggregatedData())

	r.Delete("/{uuid}", auth.JWTHandlerFunc(
		auth.RestrictHandlerFunc(s.deleteData(), permission.DataDelete),
		tokenService,
		roleService,
	))

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAggregatedData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataType := r.URL.Query().Get("type")
		meshNodeUUIDs := r.URL.Query()["meshNodes"]

		endTime := time.Now()
		measuredEnd := r.URL.Query().Get("measuredEnd")
		if measuredEnd != "" {
			var err error
			endTime, err = time.Parse(time.RFC3339, measuredEnd)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 measuredStart in wrong time format"))
				return
			}
		}

		//Falls keine Startzeit angegeben übernimmt er als Startzeit Now() vor einem Tag
		//Falls allerdings eine Endzeit gegeben ist, dann nimmt er diese vor einem Tag
		startTime := time.Now().AddDate(0, 0, -1)
		if measuredStart := r.URL.Query().Get("measuredStart"); measuredStart != "" {
			var err error
			startTime, err = time.Parse(time.RFC3339, measuredStart)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 measuredStart in wrong time format"))
				return
			}
		} else if measuredEnd != "" {
			startTime = endTime.AddDate(0, 0, -1)
		}

		sampleTime := time.Duration(0)
		if sampleDuration := r.URL.Query().Get("sampleDuration"); sampleDuration != "" {
			var err error
			sampleTime, err = time.ParseDuration(sampleDuration)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 sampleDuration in wrong time format"))
				return
			}
		}

		sampleCount := 0
		if sampleCountValue := r.URL.Query().Get("sampleCount"); sampleCountValue != "" {
			var err error
			sampleCount, err = strconv.Atoi(sampleCountValue)
			if err != nil || sampleCount <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 wrong input for sampleCount. Must be integer greater than 0."))
				return
			}
		}

		if sampleTime == 0 && sampleCount == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 eiher sampleDuration or sampleCount must be given"))
			return
		}

		if sampleCount != 0 && sampleTime != 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 sampleDuration or sampleCount are mutually exclusive"))
			return
		}

		aggregateFunction := r.URL.Query().Get("aggregateFunction")
		if !isValidAggregateFunction(aggregateFunction) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 aggregateFunction is required"))
			return
		}

		data, err := s.dataStore.GetAggregatedData(dataType, meshNodeUUIDs, startTime, endTime, sampleTime, sampleCount, aggregateFunction)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		if data.Samples == nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("404 no samles found"))
			return
		}

		response, err := json.Marshal(data)
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

func isValidAggregateFunction(aggregateFunction string) bool {
	validFunctions := map[string]bool{
		"count":   true,
		"sum":     true,
		"minimum": true,
		"maximum": true,
		"average": true,
		"range":   true,
		"median":  true,
	}

	_, ok := validFunctions[aggregateFunction]
	return ok
}

func (s service) getManyData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataType := r.URL.Query().Get("type")
		meshNodeUUIDs := r.URL.Query()["meshNodes"]
		measuredStart := r.URL.Query().Get("measuredStart")

		var startTime time.Time
		if measuredStart == "" {
			startTime = time.Unix(0, 0)
		} else {
			var err error
			startTime, err = time.Parse(time.RFC3339, measuredStart)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 measuredStart in wrong time format"))
				return
			}
		}

		var endTime time.Time
		measuredEnd := r.URL.Query().Get("measuredEnd")
		if measuredEnd == "" {
			endTime = time.Unix(0, 0)
		} else {
			var err error
			endTime, err = time.Parse(time.RFC3339, measuredEnd)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 measuredEnd in wrong time format"))
				return
			}
		}

		data, err := s.dataStore.GetManyData(dataType, meshNodeUUIDs, startTime, endTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		if data.MeasuredDatas == nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("404 no data found"))
			return
		}

		response, err := json.Marshal(data)
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

func (s service) getData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := chi.URLParam(r, "uuid")

		data, err := s.dataStore.GetData(uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		if data.Value == "" {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("404 no data found"))
			return
		}

		response, err := json.Marshal(data)
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

func (s service) getDataTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataTypes, err := s.dataStore.GetTypes()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		if dataTypes == nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("404 no dataType found"))
			return
		}

		response, err := json.Marshal(dataTypes)
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

func (s service) deleteData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := chi.URLParam(r, "uuid")

		err := s.dataStore.DeleteData(uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
