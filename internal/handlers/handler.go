package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"strconv"
)

var localStorage MemStorage = NewMemStorage()

func InitRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	r.Route("/", func(r chi.Router) {
		r.Get("/", mainHandle)
		r.Route("/update", func(r chi.Router) {
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/{metricName}", func(r chi.Router) {
					r.Route("/{metricValue}", func(r chi.Router) {
						r.Post("/", UpdateMetric)

					})
				})

			})
		})
	})
	return r
}

func mainHandle(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.Path)
}

func NewMetricHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update/", OnlyPostAllowed(http.HandlerFunc(UpdateMetric))))
	return mux
}

// MemStorage Keeps gauge and counter
type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

// NewMemStorage — constructor of the type MemStorage.
func NewMemStorage() MemStorage {
	return MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

// UpdateMetric handles update metrics request
func UpdateMetric(res http.ResponseWriter, req *http.Request) {

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")
	fmt.Println(metricType, metricName)

	switch metricType {
	case "gauge":
		valueFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			io.WriteString(res, fmt.Sprintf("ERROR: %v\n", err))
			return
		}
		localStorage.gauge[metricName] = valueFloat
		res.Header().Set("Content-type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, fmt.Sprintf("%0.2f\n", localStorage.gauge[metricName]))
	case "counter":
		valueInt, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			io.WriteString(res, fmt.Sprintf("ERROR: %v\n", err))
			return
		}
		localStorage.counter[metricName] += valueInt
		res.Header().Set("Content-type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, fmt.Sprintf("%d\n", localStorage.counter[metricName]))
	default:
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, "ERROR: UNKNOWN METRIC TYPE. Only gauge and counter are available\n")
		return
	}

}
