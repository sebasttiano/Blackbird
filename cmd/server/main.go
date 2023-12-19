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
	err := http.ListenAndServe(`localhost:8080`, handlers.InitRouter())
	return err
}
