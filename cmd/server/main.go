// Package main серверная часть, поднимает Api, принимает и отдает метрики, хранит либо в БД, либо в памяти.
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/config"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"go.uber.org/zap"
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
		if err := currentApp.service.Save(); err != nil {
			logger.Log.Error("couldn`t finally save file after graceful shutdown", zap.Error(err))
		}
	}
}

// run инициализирует заисимости и запускает http сервер.
func run(cfg *config.Config) error {

	serviceSettings := &service.ServiceSettings{SaveFilePath: cfg.FileStoragePath, Retries: cfg.RetriesDB, BackoffFactor: cfg.BackoffFactor}
	if cfg.DatabaseDSN != "" {
		var conn *sqlx.DB
		conn, err := sqlx.Connect("pgx", cfg.DatabaseDSN)
		if err != nil {
			logger.Log.Error("database openning failed", zap.Error(err))
			os.Exit(1)
		}
		defer conn.Close()
		serviceSettings.Conn = conn
		serviceSettings.DBSave = true
	} else if cfg.FileStoragePath != "" {
		serviceSettings.FileSave = true
	}

	if cfg.StoreInterval == 0 {
		serviceSettings.SyncSave = true
	}

	if err := currentApp.Initialize(serviceSettings, cfg.SecretKey); err != nil {
		logger.Log.Error("failed to init app", zap.Error(err))
	}

	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Second * time.Duration(cfg.StoreInterval))
		go service.TickerSaver(ticker, currentApp.service)
	}

	if *cfg.RestoreMetrics && serviceSettings.FileSave {
		if err := currentApp.service.Restore(); err != nil {
			logger.Log.Error("couldn`t restore data")
		}
		logger.Log.Debug("metrics were restored")
	}

	logger.Log.Info("Running server", zap.String("address", cfg.ServerIPAddr))
	return http.ListenAndServe(cfg.ServerIPAddr, currentApp.views.InitRouter())

}
