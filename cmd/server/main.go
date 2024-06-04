// Package main серверная часть, поднимает Api, принимает и отдает метрики, хранит либо в БД, либо в памяти.
package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/server"
	"net"
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
	serviceSettings := &service.ServiceSettings{SaveFilePath: cfg.FileStoragePath, Retries: cfg.RetriesDB, BackoffFactor: cfg.BackoffFactor, TrustedSubnet: nil}
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

	if cfg.TrustedSubnet != "" {
		_, subnet, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			logger.Log.Info("trusted subnet parse failed", zap.Error(err))
		} else {
			logger.Log.Info("trusted subnet parsed", zap.String("subnet", subnet.String()))
			serviceSettings.TrustedSubnet = subnet
		}
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	if cfg.GRPSServerIPAddr != "" {
		grpcSrv := server.NewGRPSServer(currentApp.service)
		wg.Add(1)
		go grpcSrv.Start(cfg.GRPSServerIPAddr)
		go grpcSrv.HandleShutdown(ctx, wg)
	}

	go srv.Start(cfg)
	go srv.HandleShutdown(ctx, wg, cfg)

	wg.Wait()
}
