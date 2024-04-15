// Package main Агент собирает с опеределенным интервалом метрики с локальной машины и пушит их с опереденным интерфвалом на сервер.
package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/sebasttiano/Blackbird.git/internal/agent"
	"github.com/sebasttiano/Blackbird.git/internal/config"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
)

func main() {

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
	go run(&cfg) // запускаем сервер

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
	logger.Log.Info(fmt.Sprintf("Metric repository server address is set to %s\n", cfg.ServerIPAddr))
	a := agent.NewAgent("http://"+cfg.ServerIPAddr, 3, 1, cfg.SecretKey)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	jobsMetrics := make(chan agent.MetricsSet, 10)
	jobsGMetrics := make(chan agent.GopsutilMetricsSet, 10)

	a.WG.Add(2)
	go a.GetMetrics(ctx, time.Duration(cfg.PollInterval)*time.Second, jobsMetrics)
	go a.GetGopsutilMetrics(ctx, time.Duration(cfg.PollInterval)*time.Second, jobsGMetrics)

	for i := 0; i < int(cfg.RateLimit); i++ {
		a.WG.Add(1)
		go a.IterateStructFieldsAndSend(ctx, time.Duration(cfg.ReportInterval)*time.Second, jobsMetrics, jobsGMetrics)
	}
	a.WG.Wait()
	cfg.WG.Done()
	return nil
}
