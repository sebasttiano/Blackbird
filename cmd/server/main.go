package main

import (
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"net/http"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, handlers.InitRouter())
}
