package main

import (
	"flag"
)

var serverIPAddr string
var pollInterval int64
var reportInterval int64

// parseFlags handles args of cli
func parseFlags() {
	flag.StringVar(&serverIPAddr, "a", "localhost:8080", "address and port of metric storage server")
	flag.Int64Var(&pollInterval, "p", 2, "interval in seconds between poll requests")
	flag.Int64Var(&reportInterval, "r", 10, "interval in seconds between push requests to server")
	flag.Parse()
}
