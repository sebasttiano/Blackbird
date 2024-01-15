package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/common"
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
	parseFlags()
	go run()

	<-done
	logger.Log.Debug("Shutdown signal interrupted")
	if flagFileStoragePath != "" {
		localStorage := *storage.GetCurrentStorage()
		if err := localStorage.SaveToFile(flagFileStoragePath); err != nil {
			logger.Log.Error("couldn`t finally save file after graceful shutdown", zap.Error(err))
		}
	}
}

func run() error {
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	localStorage := *storage.GetCurrentStorage()
	settings := storage.GetCurrentServerSettings()

	if flagStoreInterval > 0 && flagFileStoragePath != "" {
		ticker := time.NewTicker(time.Second * time.Duration(flagStoreInterval))
		go common.Schedule(ticker, flagFileStoragePath)
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
