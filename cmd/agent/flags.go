package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	serverIPAddr           string
	flagLogLevel           string
	pollInterval           int64
	reportInterval         int64
	httpClientRetry             = 3
	httpClientRetryBackoff uint = 1
)

// parseFlags handles args of cli
func parseFlags() {
	// Parse from cli
	flag.StringVar(&serverIPAddr, "a", "localhost:8080", "address and port of metric storage server")
	flag.Int64Var(&pollInterval, "p", 2, "interval in seconds between poll requests")
	flag.Int64Var(&reportInterval, "r", 10, "interval in seconds between push requests to server")
	flag.Parse()

	if envServerIPAddr := os.Getenv("ADDRESS"); envServerIPAddr != "" {
		serverIPAddr = envServerIPAddr
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollInterval, _ = strconv.ParseInt(envPollInterval, 10, 64)
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		reportInterval, _ = strconv.ParseInt(envReportInterval, 10, 64)
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
}
