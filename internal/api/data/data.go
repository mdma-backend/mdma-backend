package data

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/types/permission"
)

type DataStore interface {
	GetAggregatedData(dataType string, meshNodeUUIDs []string, measuredStart string, measuredEnd string, sampleDuration string, sampleCount int, aggregateFunction string) (AggregatedData, error)
	GetManyData(dataType string, meshNodeUUIDs []string, measuredStart string, measuredEnd string) (ManyData, error)
	GetData(uuid string) (Data, error)
	DeleteData(uuid string) error
	GetTypes() ([]string, error)
}

type service struct {
	handler   http.Handler
	dataStore DataStore
}

type AggregatedData struct {
	AggregateFunction string
	DataType          string
	MeshNodeUUIDs     []string
	Samples           []Sample
}

type Sample struct {
	FirstMeasurementAt string
	LastMeasurementAt  string
	Value              string
}

type ManyData struct {
	DataType      string
	MeasuredDatas []MeasuredData
}

type MeasuredData struct {
	MeshnodeUUID string
	Measurements []Measurement
}

type Measurement struct {
	UUID       string
	MeasuredAt string
	Value      string
}

type Data struct {
	UUID           string
	ControllerUuid string
	Type           string
	CreatedAt      string
	MeasuredAt     string
	Value          string
}

func NewService(dataStore DataStore) http.Handler {
	r := chi.NewRouter()
	s := service{
		handler:   r,
		dataStore: dataStore,
	}

	r.Get("/", auth.RestrictHandlerFunc(s.getManyData(), permission.DataRead))
	r.Get("/{uuid}", auth.RestrictHandlerFunc(s.getData(), permission.DataRead))
	r.Get("/types", auth.RestrictHandlerFunc(s.getDataTypes(), permission.DataRead))
	r.Get("/aggregated", auth.RestrictHandlerFunc(s.getAggregatedData(), permission.DataRead))
	r.Delete("/{uuid}", auth.RestrictHandlerFunc(s.deleteData(), permission.DataDelete))

	return s
}

func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s service) getAggregatedData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataType := r.URL.Query().Get("type")
		meshNodeUUIDs := r.URL.Query()["meshNodes"]

		measuredStart := r.URL.Query().Get("measuredStart")
		if measuredStart == "" {
			measuredStart = time.Unix(0, 0).String()
		}
		measuredEnd := r.URL.Query().Get("measuredEnd")
		if measuredEnd == "" {
			measuredEnd = time.Unix(0, 0).String()
		}

		sampleDuration := r.URL.Query().Get("sampleDuration")
		sampleCount := 0
		if sampleCountValue := r.URL.Query().Get("sampleCount"); sampleCountValue != "" {
			var err error
			sampleCount, err = strconv.Atoi(sampleCountValue)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("400 wrong input for sampleCount. Must be integer greater than 0."))
				return
			}
		}

		if sampleDuration == "" && sampleCount == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("400 eiher sampleDuration or sampleCount must be given"))
			return
		}

		if sampleCount != 0 && sampleDuration != "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("400 sampleDuration or sampleCount are mutually exclusive"))
			return
		}

		aggregateFunction := r.URL.Query().Get("aggregateFunction")
		if aggregateFunction == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("400 aggregateFunction is required"))
			return
		}

		println(dataType, meshNodeUUIDs, measuredStart, measuredEnd, sampleDuration, sampleCount, aggregateFunction)

		data, err := s.dataStore.GetAggregatedData(dataType, meshNodeUUIDs, measuredStart, measuredEnd, sampleDuration, sampleCount, aggregateFunction)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
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

func (s service) getManyData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataType := r.URL.Query().Get("type")
		meshNodeUUIDs := r.URL.Query()["meshNodes"]
		measuredStart := r.URL.Query().Get("measuredStart")
		if measuredStart == "" {
			measuredStart = time.Unix(0, 0).String()
		}
		measuredEnd := r.URL.Query().Get("measuredEnd")
		if measuredEnd == "" {
			measuredEnd = time.Unix(0, 0).String()
		}

		data, err := s.dataStore.GetManyData(dataType, meshNodeUUIDs, measuredStart, measuredEnd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
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

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data deleted successfully"))
	}
}
