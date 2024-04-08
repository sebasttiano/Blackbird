package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/config"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.NewServerConfig()

	if err != nil {
		fmt.Printf("config initialization failed: %v", err)
		return
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		fmt.Println("logger initialization failed")
		return
	}

	go run(cfg)

	<-done
	logger.Log.Debug("shutdown signal interrupted")
	if cfg.FileStoragePath != "" {
		if err := currentApp.store.Save(); err != nil {
			logger.Log.Error("couldn`t finally save file after graceful shutdown", zap.Error(err))
		}
	}
}

// run init dependencies and starts http server
func run(cfg config.Config) error {

	storeSettings := &repository.StoreSettings{SaveFilePath: cfg.FileStoragePath, Retries: cfg.RetriesDB, BackoffFactor: cfg.BackoffFactor}
	if cfg.DatabaseDSN != "" {
		var conn *sqlx.DB
		conn, err := sqlx.Connect("pgx", cfg.DatabaseDSN)
		if err != nil {
			logger.Log.Error("database openning failed", zap.Error(err))
			os.Exit(1)
		}
		defer conn.Close()
		storeSettings.Conn = conn
		storeSettings.DBSave = true
	} else if cfg.FileStoragePath != "" {
		storeSettings.FileSave = true
	}

	if cfg.StoreInterval == 0 {
		storeSettings.SyncSave = true
	}

	if err := currentApp.Initialize(storeSettings, cfg.SecretKey); err != nil {
		logger.Log.Error("failed to init app", zap.Error(err))
	}

	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Second * time.Duration(cfg.StoreInterval))
		go repository.TickerSaver(ticker, currentApp.store)
	}

	if *cfg.RestoreMetrics && storeSettings.FileSave {
		if err := currentApp.store.Restore(); err != nil {
			logger.Log.Error("couldn`t restore data")
		}
		logger.Log.Debug("metrics were restored")
	}

	logger.Log.Info("Running server", zap.String("address", cfg.ServerIPAddr))
	return http.ListenAndServe(cfg.ServerIPAddr, currentApp.views.InitRouter())

}
