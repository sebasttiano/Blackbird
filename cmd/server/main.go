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

type Middleware func(http.Handler) http.Handler

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

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Method: %s\r\n", req.Method)
	body += "Header ===============\r\n"
	for k, v := range req.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "Query parameters ===============\r\n"
	for k, v := range req.URL.Query() {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	res.Write([]byte(body))
}

func updateGauge(res http.ResponseWriter, req *http.Request) {

	incomingParams := strings.Split(req.URL.Path, "/")
	if len(incomingParams) != 2 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	metricName := incomingParams[0]
	valueFloat, err := strconv.ParseFloat(incomingParams[1], 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, fmt.Sprintf("ERROR: %v\n", err))
		return
	}

	localStorage.gauge[metricName] = valueFloat
	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusOK)
	io.WriteString(res, fmt.Sprintf("%f\n", localStorage.gauge[metricName]))
}

func updateCounter(res http.ResponseWriter, req *http.Request) {

	incomingParams := strings.Split(req.URL.Path, "/")
	if len(incomingParams) != 2 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	metricName := incomingParams[0]
	valueInt, err := strconv.ParseInt(incomingParams[1], 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		io.WriteString(res, fmt.Sprintf("ERROR: %v\n", err))
		return
	}

	localStorage.counter[metricName] += valueInt
	res.Header().Set("Content-type", "text/plain")
	res.WriteHeader(http.StatusOK)
	io.WriteString(res, fmt.Sprintf("%d\n", localStorage.counter[metricName]))
}

func updatePage(res http.ResponseWriter, req *http.Request) {
	var body string

	body += "ANOTHER SECTION ===============\r\n"
	for k, v := range req.URL.Query() {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}

	res.Write([]byte(body))
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	mux := http.NewServeMux()
	mux.Handle("/", middleware(http.HandlerFunc(mainPage)))
	mux.Handle("/update/", middleware(http.HandlerFunc(updatePage)))
	mux.Handle("/update/gauge/", http.StripPrefix("/update/gauge/", middleware(http.HandlerFunc(updateGauge))))
	mux.Handle("/update/counter/", http.StripPrefix("/update/counter/", middleware(http.HandlerFunc(updateCounter))))

	err := http.ListenAndServe(`localhost:8080`, mux)
	return err
}
