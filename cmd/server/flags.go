package main

import (
	"flag"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"os"
	"strconv"
)

// non-export var flagRunAddr keeps ip address and port to run server on
// flagLogLevel keeps level of the logger
var (
	flagRunAddr         string
	flagLogLevel        string
	flagFileStoragePath string
	flagStoreInterval   int
	flagRestoreOnStart  bool
	flagDatabaseDSN     string
)

// parseFlags handles args of cli
func parseFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagStoreInterval, "i", 300, "set interval in seconds to write metrics in file")
	flag.StringVar(&flagFileStoragePath, "f", "/tmp/metrics-db.json", "specify the file to save metrics to")
	flag.BoolVar(&flagRestoreOnStart, "r", true, "Restore saved metrics on start")
	flag.StringVar(&flagDatabaseDSN, "d", "", "database host connect to, user and password")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagFileStoragePath = envFileStoragePath
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		var err error
		flagStoreInterval, err = strconv.Atoi(envStoreInterval)
		if err != nil {
			logger.Log.Error("you must pass int to env var store_interval")
		}
	}

	if envRestoreOnStart := os.Getenv("RESTORE"); envRestoreOnStart != "" {
		var err error
		flagRestoreOnStart, err = strconv.ParseBool(envRestoreOnStart)
		if err != nil {
			logger.Log.Error("you must pass bool value to env var restore")
		}
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		flagDatabaseDSN = envDatabaseDSN
	}
	return nil
}
