package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		panic(err)
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
