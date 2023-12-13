package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var localStorage MemStorage = NewMemStorage()

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

// middleware gets Handler, makes some validation and returns also Handler.
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Only POST method allowed
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(res, req)
	})
}

// updateMetric handles update metrics request
func updateMetric(res http.ResponseWriter, req *http.Request) {

	incomingParams := strings.Split(req.URL.Path, "/")
	if len(incomingParams) != 3 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	switch incomingParams[0] {
	case "gauge":
		metricName := incomingParams[1]
		valueFloat, err := strconv.ParseFloat(incomingParams[2], 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			io.WriteString(res, fmt.Sprintf("ERROR: %v\n", err))
			return
		}
		localStorage.gauge[metricName] = valueFloat
		res.Header().Set("Content-type", "text/plain")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, fmt.Sprintf("%f\n", localStorage.gauge[metricName]))
	case "counter":
		metricName := incomingParams[1]
		valueInt, err := strconv.ParseInt(incomingParams[2], 10, 64)
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
func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update/", middleware(http.HandlerFunc(updateMetric))))

	err := http.ListenAndServe(`localhost:8080`, mux)
	return err
}
