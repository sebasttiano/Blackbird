// Package config парсит переменные окружения и флаги при запуске приложения. Приоритет у переменных окружения.
// Генерирутеся объекты конфига для агента и сервера.
package config

import (
	"encoding/json"
	"flag"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"

	"github.com/caarlos0/env/v6"
)

// Config содержит все передаваемые переменные нужные для приложения
type Config struct {
	ServerIPAddr    string `env:"ADDRESS" json:"address"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	SecretKey       string `env:"KEY"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	RestoreMetrics  *bool  `env:"RESTORE" json:"restore"`
	PollInterval    int64  `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval  int64  `env:"REPORT_INTERVAL" json:"report_interval"`
	RateLimit       uint64 `env:"RATE_LIMIT"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile      string `env:"CONFIG"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	RetriesDB       uint
	BackoffFactor   uint
	Profiler        *bool `env:"PROFILER"`
	WG              sync.WaitGroup
}

func (c *Config) SetDefault() {

	if c.ServerIPAddr == "" {
		c.ServerIPAddr = "localhost:8080"
	}

	if c.PollInterval == 0 {
		c.PollInterval = 2
	}

	if c.ReportInterval == 0 {
		c.ReportInterval = 5
	}

	if c.RateLimit == 0 {
		c.RateLimit = 1
	}

	if c.Profiler == nil {
		f, _ := strconv.ParseBool("false")
		c.Profiler = &f
	}

	if c.StoreInterval == 0 {
		c.StoreInterval = 300
	}

	if c.RestoreMetrics == nil {
		f, _ := strconv.ParseBool("true")
		c.RestoreMetrics = &f

	}

	if c.FileStoragePath == "" {
		c.FileStoragePath = "/tmp/metrics-db.json"
	}
}

// NewAgentConfig конструктор для Config
func NewAgentConfig() (*Config, error) {
	flags := parseAgentFlags()
	config := Config{}
	configJSON := Config{}

	if err := env.Parse(&config); err != nil {
		return &Config{}, err
	}

	if config.ConfigFile == "" {
		config.ConfigFile = flags.ConfigFile
	}

	if config.ConfigFile != "" {
		data, err := os.ReadFile(config.ConfigFile)
		if err != nil {
			logger.Log.Error("failed to read config file", zap.String("file", config.ConfigFile), zap.Error(err))
			return nil, err
		}
		if err := json.Unmarshal(data, &configJSON); err != nil {
			logger.Log.Error("failed to unmarshal config file. check your json", zap.String("file", config.ConfigFile), zap.Error(err))
			return nil, err
		}
	}

	if config.ServerIPAddr == "" {
		config.ServerIPAddr = flags.ServerIPAddr
		if config.ServerIPAddr == "" {
			config.ServerIPAddr = configJSON.ServerIPAddr
		}
	}

	if config.PollInterval == 0 {
		config.PollInterval = flags.PollInterval
		if config.PollInterval == 0 {
			config.PollInterval = configJSON.PollInterval
		}
	}

	if config.ReportInterval == 0 {
		config.ReportInterval = flags.ReportInterval
		if config.ReportInterval == 0 {
			config.ReportInterval = configJSON.ReportInterval
		}
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

	if config.Profiler == nil {
		config.Profiler = flags.Profiler
	}

	if config.CryptoKey == "" {
		config.CryptoKey = flags.CryptoKey
		if config.CryptoKey == "" {
			config.CryptoKey = configJSON.CryptoKey
		}
	}

	config.SetDefault()
	return &config, nil
}

// parseAgentFlags считывает переменные с консоли для клиента
func parseAgentFlags() Config {
	// Parse from cli
	serverIPAddr := flag.String("a", "", "address and port of metric repository server")
	pollInterval := flag.Int64("p", 0, "interval in seconds between poll requests")
	reportInterval := flag.Int64("r", 0, "interval in seconds between push requests to server")
	flagSecretKey := flag.String("k", "", "secret key for digital signature")
	flagRateLimit := flag.Uint64("l", 0, "number of simultaneous requests to server")
	flagProfiler := flag.Bool("profiler", false, "enable profiler")
	flagCryptoKey := flag.String("crypto-key", "", "path to file with public key")
	flagConfigFile := flag.String("config", "", "path to config file")

	flag.Parse()

	return Config{
		ServerIPAddr:   *serverIPAddr,
		PollInterval:   *pollInterval,
		ReportInterval: *reportInterval,
		SecretKey:      *flagSecretKey,
		RateLimit:      *flagRateLimit,
		Profiler:       flagProfiler,
		CryptoKey:      *flagCryptoKey,
		ConfigFile:     *flagConfigFile,
	}
}

// NewServerConfig конструктор конфига для серверной части
func NewServerConfig() (*Config, error) {
	flags := parseServerFlags()
	config := Config{RetriesDB: 1, BackoffFactor: 1}
	configJSON := Config{}

	if err := env.Parse(&config); err != nil {
		return &Config{}, err
	}

	if config.ConfigFile == "" {
		config.ConfigFile = flags.ConfigFile
	}

	if config.ConfigFile != "" {
		data, err := os.ReadFile(config.ConfigFile)
		if err != nil {
			logger.Log.Error("failed to read config file", zap.String("file", config.ConfigFile), zap.Error(err))
			return nil, err
		}
		if err := json.Unmarshal(data, &configJSON); err != nil {
			logger.Log.Error("failed to unmarshal config file. check your json", zap.String("file", config.ConfigFile), zap.Error(err))
			return nil, err
		}
	}

	if config.ServerIPAddr == "" {
		config.ServerIPAddr = flags.ServerIPAddr
		if config.ServerIPAddr == "" {
			config.ServerIPAddr = configJSON.ServerIPAddr
		}
	}

	if config.StoreInterval == 0 {
		config.StoreInterval = flags.StoreInterval
		if config.StoreInterval == 0 {
			config.StoreInterval = configJSON.StoreInterval
		}
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = flags.FileStoragePath
		if config.FileStoragePath == "" {
			config.FileStoragePath = configJSON.FileStoragePath
		}
	}

	if config.RestoreMetrics == nil {
		config.RestoreMetrics = flags.RestoreMetrics
		if config.RestoreMetrics == nil {
			config.RestoreMetrics = configJSON.RestoreMetrics
		}
	}

	if config.DatabaseDSN == "" {
		config.DatabaseDSN = flags.DatabaseDSN
		if config.DatabaseDSN == "" {
			config.DatabaseDSN = configJSON.DatabaseDSN
		}
	}

	if config.SecretKey == "" {
		config.SecretKey = flags.SecretKey
	}

	if config.CryptoKey == "" {
		config.CryptoKey = flags.CryptoKey
		if config.CryptoKey == "" {
			config.CryptoKey = configJSON.CryptoKey
		}
	}

	if config.ConfigFile == "" {
		config.ConfigFile = flags.ConfigFile
	}

	if config.TrustedSubnet == "" {
		config.TrustedSubnet = flags.TrustedSubnet
		if config.TrustedSubnet == "" {
			config.TrustedSubnet = configJSON.TrustedSubnet
		}

	}

	config.SetDefault()
	return &config, nil
}

// parseServerFlags считывает переменные с консоли для сервера
func parseServerFlags() Config {
	// Parse from cli
	serverIPAddr := flag.String("a", "", "address and port to run server")
	serverStoreInterval := flag.Int("i", 0, "set interval in seconds to write metrics in file")
	fileStoragePath := flag.String("f", "", "specify the file to save metrics to")
	databaseDSN := flag.String("d", "", "database host connect to, user and password")
	secretKey := flag.String("k", "", "secret key for digital signature")
	cryptoKey := flag.String("crypto-key", "", "path to file with private key")
	configFile := flag.String("config", "", "path to config file")
	trustedSubnet := flag.String("t", "", "trusted subnet")

	var restoreOnStart *bool
	flag.BoolFunc("r", "restore saved metrics on start", func(restore string) error {
		if restore == "" {
			restoreOnStart = nil
		}
		value, err := strconv.ParseBool(restore)
		if err != nil {
			logger.Log.Error("failed to parse profiler flag", zap.String("restore", restore), zap.Error(err))
			restoreOnStart = nil
		}
		restoreOnStart = &value
		return nil
	})
	flag.Parse()

	return Config{
		ServerIPAddr:    *serverIPAddr,
		StoreInterval:   *serverStoreInterval,
		FileStoragePath: *fileStoragePath,
		RestoreMetrics:  restoreOnStart,
		DatabaseDSN:     *databaseDSN,
		SecretKey:       *secretKey,
		CryptoKey:       *cryptoKey,
		ConfigFile:      *configFile,
		TrustedSubnet:   *trustedSubnet,
	}
}
