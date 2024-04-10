package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"github.com/sebasttiano/Blackbird.git/templates"
	"go.uber.org/zap"
)

type ServerViews struct {
	Service   *service.Service
	templates templates.HTMLTemplates
	DB        *sqlx.DB
	SignKey   string
}

func NewServerViews(service *service.Service) ServerViews {
	return ServerViews{Service: service, templates: templates.ParseTemplates()}
}

func (s *ServerViews) InitRouter() chi.Router {

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(WithLogging, CheckSign(s.SignKey), GzipMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Get("/", s.MainHandle)
		r.Get("/ping", s.PingDB)
		r.Post("/updates/", s.UpdateMetricsJSON)
		r.Route("/value", func(r chi.Router) {
			r.Post("/", s.GetMetricJSON)
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/{metricName}", func(r chi.Router) {
					r.Get("/", s.GetMetric)
				})
			})
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", s.UpdateMetricJSON)
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/{metricName}", func(r chi.Router) {
					r.Route("/{metricValue}", func(r chi.Router) {
						r.Post("/", s.UpdateMetric)
					})
				})
			})
		})
	})
	return r
}

// MainHandle render html with all available metrics at the moment
func (s *ServerViews) MainHandle(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()

	res.Header().Set("Content-Type", "text/html")
	data := s.Service.GetAllValues(ctx)
	if err := s.templates.IndexTemplate.Execute(res, data); err != nil {
		logger.Log.Error("couldn`t render the html template", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// GetMetric gets metric from repository via interface method and sends in a
// response
func (s *ServerViews) GetMetric(res http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	value, err := s.Service.GetValue(ctx, metricName, metricType)
	if err != nil {
		logger.Log.Error("couldn`t find requested metric. ", zap.Error(err))
		http.Error(res, err.Error(), http.StatusNotFound)
	} else {

		if s.SignKey != "" {
			res.Header().Add("HashSHA256", sign(value, s.SignKey))
		}

		io.WriteString(res, fmt.Sprintf("%v\n", value))
	}
}

// GetMetricJSON gets metric from repository via interface method and sends in a model
// response
func (s *ServerViews) GetMetricJSON(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	if req.Header.Get("Content-Type") != "application/json" {
		logger.Log.Error("got request with wrong header", zap.String("Content-Type", req.Header.Get("Content-Type")))
		http.Error(res, "error: check your header Content-Type\n", http.StatusBadRequest)
	}

	logger.Log.Debug("decoding incoming request")
	var metrics models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()

	if err := s.Service.GetModelValue(ctx, &metrics); err != nil {
		logger.Log.Debug("couldn`t get model", zap.Error(err))
		http.Error(res, "couldn`t get model", http.StatusNotFound)
	}

	if s.SignKey != "" {
		res.Header().Add("HashSHA256", sign(metrics, s.SignKey))
	}

	enc := json.NewEncoder(res)
	if err := enc.Encode(metrics); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateMetric handles update metrics request
func (s *ServerViews) UpdateMetric(res http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	if err := s.Service.SetValue(ctx, metricName, metricType, metricValue); err != nil {
		logger.Log.Error("couldn`t save metric. error: ", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

// UpdateMetricJSON handles update metrics request in json format
func (s *ServerViews) UpdateMetricJSON(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	if req.Header.Get("Content-Type") != "application/json" {
		logger.Log.Error("got request with wrong header", zap.String("Content-Type", req.Header.Get("Content-Type")))
		http.Error(res, "error: check your header Content-Type", http.StatusBadRequest)
	}

	logger.Log.Debug("decoding incoming request")
	var metrics models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()

	if err := s.Service.SetModelValue(ctx, []*models.Metrics{&metrics}); err != nil {
		logger.Log.Error("couldn`t save metric. error: ", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	enc := json.NewEncoder(res)
	if err := enc.Encode(metrics); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateMetricsJSON handles batch with metrics
func (s *ServerViews) UpdateMetricsJSON(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if req.Header.Get("Content-Type") != "application/json" {
		logger.Log.Error("got request with wrong header", zap.String("Content-Type", req.Header.Get("Content-Type")))
		http.Error(res, "error: check your header Content-Type", http.StatusBadRequest)
	}
	var metrics []*models.Metrics
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()

	if err := s.Service.SetModelValue(ctx, metrics); err != nil {
		logger.Log.Error("couldn`t save metric. error: ", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	enc := json.NewEncoder(res)
	if err := enc.Encode(metrics); err != nil {
		logger.Log.Error("error encoding response", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// PingDB checks connection to database
func (s *ServerViews) PingDB(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	if err := s.DB.PingContext(ctx); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// sign any string like object with hmac signature
func sign(value any, key string) string {
	b, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write(b)
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}
