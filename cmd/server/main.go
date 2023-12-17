package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	mux := http.NewServeMux()
	mux.Handle("/update/", http.StripPrefix("/update/", handlers.Middleware(http.HandlerFunc(handlers.UpdateMetric))))

	err := http.ListenAndServe(`localhost:8080`, mux)
	return err
}
