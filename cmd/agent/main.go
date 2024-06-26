// Package main Агент собирает с опеределенным интервалом метрики с локальной машины и пушит их с опереденным интерфвалом на сервер.
package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"go.uber.org/zap"

	"github.com/sebasttiano/Blackbird.git/internal/agent"
	"github.com/sebasttiano/Blackbird.git/internal/config"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

type templateInfoEntry struct {
	Version string
	Date    string
	Commit  string
}

//go:embed agent_info.txt
var agentInfo string

func main() {
	tmpl, err := template.New("info").Parse(agentInfo)
	if err != nil {
		fmt.Printf("failed to render banner: %v", err)
	}
	tmpl.Execute(os.Stdout, templateInfoEntry{buildVersion, buildDate, buildCommit})

	cfg, err := config.NewAgentConfig()
	if err != nil {
		fmt.Printf("config initialization failed: %v", err)
		return
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		fmt.Println("logger initialization failed")
		return
	}

	cfg.WG.Add(1)
	go run(cfg)

	if *cfg.Profiler {
		if err := http.ListenAndServe("localhost:8085", nil); err != nil {
			logger.Log.Error("profiler http server initialization failed", zap.Error(err))
		}
	}
	cfg.WG.Wait()
}

// run запускает агента.
func run(cfg *config.Config) error {
	logger.Log.Info(fmt.Sprintf("Running agent with poll interval %d and report interval %d\n", cfg.PollInterval, cfg.ReportInterval))
	var srvAddr string
	if cfg.GRPSServerIPAddr != "" {
		srvAddr = cfg.GRPSServerIPAddr
	} else {
		srvAddr = cfg.ServerIPAddr
	}
	logger.Log.Info(fmt.Sprintf("Metric repository server address is set to %s\n", srvAddr))

	var publicKey []byte
	var err error
	if cfg.CryptoKey != "" {
		publicKey, err = os.ReadFile(cfg.CryptoKey)
		if err != nil {
			logger.Log.Error("failed to read crypto key", zap.Error(err))
		}
	}

	a, err := agent.NewAgent("http://"+cfg.ServerIPAddr, 3, 1, cfg.SecretKey, publicKey, cfg.GRPSServerIPAddr)
	if err != nil && errors.Is(agent.ErrInitSender, err) {
		logger.Log.Error("failed to initialize agent", zap.Error(err))
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	jobsMetrics := make(chan agent.MetricsSet, 10)
	jobsGMetrics := make(chan agent.GopsutilMetricsSet, 10)

	a.WG.Add(2)
	go a.GetMetrics(ctx, time.Duration(cfg.PollInterval)*time.Second, jobsMetrics)
	go a.GetGopsutilMetrics(ctx, time.Duration(cfg.PollInterval)*time.Second, jobsGMetrics)

	for i := 0; i < int(cfg.RateLimit); i++ {
		a.WG.Add(1)
		go a.SendMetrics(ctx, time.Duration(cfg.ReportInterval)*time.Second, jobsMetrics, jobsGMetrics)
	}
	a.WG.Wait()
	cfg.WG.Done()
	return nil
}
