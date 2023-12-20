package main

import (
	"flag"
)

// non-export var flagRunAddr keeps ip address and port to run server on
var flagRunAddr string

// parseFlags handles args of cli
func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}
