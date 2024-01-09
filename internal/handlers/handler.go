package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"github.com/sebasttiano/Blackbird.git/templates"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type ServerFacility struct {
	localStorage  storage.HandleMemStorage
	htmlTemplates templates.HTMLTemplates
}

func NewServerFacility() ServerFacility {
	return ServerFacility{
		localStorage: &storage.MemStorage{
			Gauge:   make(map[string]float64),
			Counter: make(map[string]int64),
		},
		htmlTemplates: templates.ParseTemplates()}
}

var SrvFacility = NewServerFacility()

// InitRouter provides url and method schema and returns chi.Router
func InitRouter() chi.Router {

	r := chi.NewRouter()

	r.Use(middleware.RealIP)

	r.Route("/", func(r chi.Router) {
		r.Get("/", MainHandle)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", GetMetricJSON)
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/{metricName}", func(r chi.Router) {
					r.Get("/", GetMetric)
				})
			})
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", UpdateMetricJSON)
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/metricName", func(r chi.Router) {
					r.Post("/", UpdateMetric)
				})
			})
		})
	})
	return r
}

// MainHandle render html with all available metrics at the moment
func MainHandle(res http.ResponseWriter, req *http.Request) {

	if err := SrvFacility.htmlTemplates.IndexTemplate.Execute(res, SrvFacility.localStorage); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// GetMetric gets metric from storage via interface method and sends in a
// response
func GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	value, err := SrvFacility.localStorage.GetValue(metricName, metricType)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}

	res.WriteHeader(http.StatusOK)
	io.WriteString(res, fmt.Sprintf("%v\n", value))

}

// GetMetricJSON gets metric from storage via interface method and sends in a model
// response
func GetMetricJSON(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		logger.Log.Error("got request with wrong header", zap.String("Content-Type", req.Header.Get("Content-Type")))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log.Debug("decoding incoming request")
	var metrics models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := SrvFacility.localStorage.GetModelValue(&metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(res)
	if err := enc.Encode(metrics); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateMetric handles update metrics request
func UpdateMetric(res http.ResponseWriter, req *http.Request) {

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	if err := SrvFacility.localStorage.SetValue(metricName, metricType, metricValue); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// UpdateMetricJSON handles update metrics request in json format
func UpdateMetricJSON(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-Type") != "application/json" {
		logger.Log.Error("got request with wrong header", zap.String("Content-Type", req.Header.Get("Content-Type")))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log.Debug("decoding incoming request")
	var metrics models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := SrvFacility.localStorage.SetModelValue(&metrics); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(res)
	if err := enc.Encode(metrics); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
