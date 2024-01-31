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
	if err := parseFlags(); err != nil {
		fmt.Printf("parsing flags failed: %s\n", err.Error())
	}
	if err := logger.Initialize(flagLogLevel); err != nil {
		fmt.Println("logger initialization failed")
		return
	}

	go run()

	<-done
	logger.Log.Debug("shutdown signal interrupted")
	if flagFileStoragePath != "" {
		if err := currentApp.store.Save(); err != nil {
			logger.Log.Error("couldn`t finally save file after graceful shutdown", zap.Error(err))
		}
	}
}

// run init dependencies and starts http server
func run() error {

	storeSettings := &storage.StoreSettings{SaveFilePath: flagFileStoragePath}

	var conn *sql.DB
	conn, err := sql.Open("pgx", flagDatabaseDSN)
	if err != nil {
		logger.Log.Error("database openning failed", zap.Error(err))
	}
	defer conn.Close()
	storeSettings.Conn = conn

	if flagDatabaseDSN != "" {
		storeSettings.DBSave = true
	} else if flagFileStoragePath != "" {
		storeSettings.FileSave = true
	}

	if flagStoreInterval == 0 {
		storeSettings.SyncSave = true
	}

	currentApp.Initialize(storeSettings)

	if flagStoreInterval > 0 {
		ticker := time.NewTicker(time.Second * time.Duration(flagStoreInterval))
		go storage.TickerSaver(ticker, currentApp.store)
	}

	if flagRestoreOnStart && storeSettings.FileSave {
		if err := currentApp.store.Restore(); err != nil {
			logger.Log.Error("couldn`t restore data")
		}
		logger.Log.Debug("metrics were restored")
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	return http.ListenAndServe(flagRunAddr, handlers.WithLogging(handlers.GzipMiddleware(currentApp.views.InitRouter())))
}
