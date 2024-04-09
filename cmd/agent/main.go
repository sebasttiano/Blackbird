package main

import (
	"context"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/agent"
	"github.com/sebasttiano/Blackbird.git/internal/config"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"
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

	if err := run(cfg); err != nil {
		logger.Log.Error("While executing agent, error occurred", zap.Error(err))
	}

}

func run(cfg config.Config) error {

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
	return nil
}
