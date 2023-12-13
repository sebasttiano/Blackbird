package main

import (
	"fmt"
	"net/http"
)

var localStorage MemStorage

// MemStorage Keeps gauge and counter
type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
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
	res.Write([]byte(req.URL.Path + "\n"))
}

func updateCounter(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(req.URL.Path + "\n"))
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
	mux.Handle("/update/gauge/", middleware(http.HandlerFunc(updateGauge)))
	mux.Handle("/update/counter/", middleware(http.HandlerFunc(updateCounter)))

	err := http.ListenAndServe(`localhost:8080`, mux)
	return err
}
