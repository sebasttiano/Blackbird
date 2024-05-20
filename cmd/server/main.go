// Package main серверная часть, поднимает Api, принимает и отдает метрики, хранит либо в БД, либо в памяти.
package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/server"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/template"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/config"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"go.uber.org/zap"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

type templateInfoEntry struct {
	Version string
	Date    string
	Commit  string
}

//go:embed server_info.txt
var serverInfo string

func main() {
	tmpl, err := template.New("info").Parse(serverInfo)
	if err != nil {
		fmt.Printf("failed to render banner: %v", err)
	}
	tmpl.Execute(os.Stdout, templateInfoEntry{buildVersion, buildDate, buildCommit})
	cfg, err := config.NewServerConfig()

	if err != nil {
		fmt.Printf("config initialization failed: %v", err)
		return
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		fmt.Println("logger initialization failed")
		return
	}

	run(cfg)
}

// run инициализирует заисимости и запускает http сервер.
func run(cfg *config.Config) {
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

	var privateKey []byte
	var err error
	if cfg.CryptoKey != "" {
		privateKey, err = os.ReadFile(cfg.CryptoKey)
		if err != nil {
			logger.Log.Error("failed to read crypto key", zap.Error(err))
		}
	}

	if err := currentApp.Initialize(serviceSettings, cfg.SecretKey, privateKey); err != nil {
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

	srv := server.NewServer(cfg.ServerIPAddr, &currentApp.views, currentApp.views.InitRouter())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go srv.Start(cfg)
	go srv.HandleShutdown(ctx, wg, cfg)

	wg.Wait()
	//logger.Log.Info("Running server", zap.String("address", cfg.ServerIPAddr))
	//return http.ListenAndServe(cfg.ServerIPAddr, currentApp.views.InitRouter())
}
