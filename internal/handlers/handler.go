package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"html/template"
	"io"
	"net/http"
	"path"
)

var localStorage = storage.NewMemStorage()

// InitRouter provides url and method schema and returns chi.Router
func InitRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	r.Route("/", func(r chi.Router) {
		r.Get("/", MainHandle)
		r.Route("/value", func(r chi.Router) {
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/{metricName}", func(r chi.Router) {
					r.Get("/", GetMetric)
				})
			})
		})
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

// MainHandle render html with all available metrics at the moment
func MainHandle(res http.ResponseWriter, req *http.Request) {

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(res, localStorage); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// GetMetric gets metric from storage via interface method and sends in a
// response
func GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	value, err := localStorage.GetValue(metricName, metricType)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}

	res.WriteHeader(http.StatusOK)
	io.WriteString(res, fmt.Sprintf("%v\n", value))

}

// NewMetricHandler custom handler-mux
func NewMetricHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update/", OnlyPostAllowed(http.HandlerFunc(UpdateMetric))))
	return mux
}

// UpdateMetric handles update metrics request
func UpdateMetric(res http.ResponseWriter, req *http.Request) {

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	if err := localStorage.SetValue(metricName, metricType, metricValue); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}
