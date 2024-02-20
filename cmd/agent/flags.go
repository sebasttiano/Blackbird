package main

import (
	"flag"
	"os"
	"strconv"
)

//var (
//	serverIPAddr           string
//	flagLogLevel           string
//	pollInterval           int64
//	reportInterval         int64
//	flagSecretKey          string
//	flagRateLimit          uint64
//	httpClientRetry             = 3
//	httpClientRetryBackoff uint = 1
//)

type Config struct {
	serverIPAddr   string
	flagLogLevel   string
	pollInterval   int64
	reportInterval int64
	flagSecretKey  string
	flagRateLimit  uint64
}

func NewConfig() Config {

	config := parseFlags()

	if envServerIPAddr := os.Getenv("ADDRESS"); envServerIPAddr != "" {
		config.serverIPAddr = envServerIPAddr
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		config.reportInterval, _ = strconv.ParseInt(envPollInterval, 10, 64)
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		config.reportInterval, _ = strconv.ParseInt(envReportInterval, 10, 64)
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		config.flagLogLevel = envLogLevel
	}

	if envSecretKey := os.Getenv("KEY"); envSecretKey != "" {
		config.flagSecretKey = envSecretKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		config.flagRateLimit, _ = strconv.ParseUint(envRateLimit, 10, 64)
	}
	return config
}

// parseFlags handles args of cli
func parseFlags() Config {
	// Parse from cli
	serverIPAddr := flag.String("a", "localhost:8080", "address and port of metric storage server")
	pollInterval := flag.Int64("p", 2, "interval in seconds between poll requests")
	reportInterval := flag.Int64("r", 5, "interval in seconds between push requests to server")
	flagSecretKey := flag.String("k", "", "secret key for digital signature")
	flagRateLimit := flag.Uint64("l", 1, "number of simultaneous requests to server")

	flag.Parse()

	return Config{
		serverIPAddr:   *serverIPAddr,
		pollInterval:   *pollInterval,
		reportInterval: *reportInterval,
		flagSecretKey:  *flagSecretKey,
		flagRateLimit:  *flagRateLimit,
	}
}
