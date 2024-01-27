package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
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
	if err := logger.Initialize(flagLogLevel); err != nil {
		fmt.Println("logger initialization failed")
		return
	}
	if err := parseFlags(); err != nil {
		logger.Log.Error("parsing flags failed: ", zap.Error(err))
	}

	var err error
	storage.DB, err = sql.Open("pgx", flagDatabaseDSN)
	defer storage.DB.Close()

	if err != nil {
		logger.Log.Error("database openning failed", zap.Error(err))
	}

	go run()

	<-done
	logger.Log.Debug("shutdown signal interrupted")
	if flagFileStoragePath != "" {
		localStorage := storage.SrvFacility.LocalStorage
		if err := localStorage.SaveToFile(flagFileStoragePath); err != nil {
			logger.Log.Error("couldn`t finally save file after graceful shutdown", zap.Error(err))
		}
	}
}

func run() error {
	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	localStorage := storage.SrvFacility.LocalStorage
	settings := storage.GetCurrentServerSettings()

	if flagStoreInterval > 0 && flagFileStoragePath != "" {
		ticker := time.NewTicker(time.Second * time.Duration(flagStoreInterval))
		go storage.TickerSaver(ticker, flagFileStoragePath)
	}

	if flagFileStoragePath == "" {
		settings.SyncSave = false
	}
	if flagStoreInterval == 0 {
		settings.SyncSave = true
		settings.SaveFilePath = flagFileStoragePath
	}
	if flagRestoreOnStart && flagFileStoragePath != "" {
		if err := localStorage.RestoreFromFile(flagFileStoragePath); err != nil {
			logger.Log.Error("couldn`t restore data from file")
		}
		logger.Log.Debug("metrics were restored from the file")
	}

	return http.ListenAndServe(flagRunAddr, handlers.WithLogging(handlers.GzipMiddleware(handlers.InitRouter())))
}
