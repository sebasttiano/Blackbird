package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerIPAddr    string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	SecretKey       string `env:"KEY"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	RestoreMetrics  *bool  `env:"RESTORE"`
	PollInterval    int64  `env:"POLL_INTERVAL"`
	ReportInterval  int64  `env:"REPORT_INTERVAL"`
	RateLimit       uint64 `env:"RATE_LIMIT"`
	RetriesDB       uint
	BackoffFactor   uint
}

func NewAgentConfig() (Config, error) {

	flags := parseAgentFlags()
	config := Config{}

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if config.ServerIPAddr == "" {
		config.ServerIPAddr = flags.ServerIPAddr
	}

	if config.PollInterval == 0 {
		config.PollInterval = flags.PollInterval
	}

	if config.ReportInterval == 0 {
		config.ReportInterval = flags.ReportInterval
	}

	if config.LogLevel == "" {
		config.LogLevel = flags.LogLevel
	}

	if config.SecretKey == "" {
		config.SecretKey = flags.SecretKey
	}

	if config.RateLimit == 0 {
		config.RateLimit = flags.RateLimit
	}
	return config, nil
}

// parseAgentFlags handles args of cli agent
func parseAgentFlags() Config {
	// Parse from cli
	serverIPAddr := flag.String("a", "localhost:8080", "address and port of metric repository server")
	pollInterval := flag.Int64("p", 2, "interval in seconds between poll requests")
	reportInterval := flag.Int64("r", 5, "interval in seconds between push requests to server")
	flagSecretKey := flag.String("k", "", "secret key for digital signature")
	flagRateLimit := flag.Uint64("l", 1, "number of simultaneous requests to server")

	flag.Parse()

	return Config{
		ServerIPAddr:   *serverIPAddr,
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
		SecretKey:      *flagSecretKey,
		RateLimit:      *flagRateLimit,
	}
}

// NewServerConfig config constructor for server
func NewServerConfig() (Config, error) {

	flags := parseServerFlags()
	config := Config{RetriesDB: 1, BackoffFactor: 1}

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if config.ServerIPAddr == "" {
		config.ServerIPAddr = flags.ServerIPAddr
	}

	if config.StoreInterval == 0 {
		config.StoreInterval = flags.StoreInterval
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = flags.FileStoragePath
	}

	if config.RestoreMetrics == nil {
		config.RestoreMetrics = flags.RestoreMetrics
	}

	if config.DatabaseDSN == "" {
		config.DatabaseDSN = flags.DatabaseDSN
	}

	if config.SecretKey == "" {
		config.SecretKey = flags.SecretKey
	}

	return config, nil
}

// parseServerFlags handles args cli for server
func parseServerFlags() Config {
	// Parse from cli
	serverIPAddr := flag.String("a", "localhost:8080", "address and port to run server")
	serverStoreInterval := flag.Int("i", 300, "set interval in seconds to write metrics in file")
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "specify the file to save metrics to")
	restoreOnStart := flag.Bool("r", true, "restore saved metrics on start")
	databaseDSN := flag.String("d", "", "database host connect to, user and password")
	secretKey := flag.String("k", "", "secret key for digital signature")

	flag.Parse()

	return Config{
		ServerIPAddr:    *serverIPAddr,
		StoreInterval:   *serverStoreInterval,
		FileStoragePath: *fileStoragePath,
		RestoreMetrics:  restoreOnStart,
		DatabaseDSN:     *databaseDSN,
		SecretKey:       *secretKey,
	}
}
