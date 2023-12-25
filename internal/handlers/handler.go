package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"github.com/sebasttiano/Blackbird.git/templates"
	"io"
	"net/http"
)

type ServerFacility struct {
	localStorage  storage.MemStorage
	htmlTemplates templates.HtmlTemplates
}

func NewServerFacility() ServerFacility {
	return ServerFacility{
		localStorage:  storage.NewMemStorage(),
		htmlTemplates: templates.ParseTemplates()}
}

var SrvFacility = NewServerFacility()

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

// UpdateMetric handles update metrics request
func UpdateMetric(res http.ResponseWriter, req *http.Request) {

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	if err := SrvFacility.localStorage.SetValue(metricName, metricType, metricValue); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.WriteHeader(http.StatusOK)
}
