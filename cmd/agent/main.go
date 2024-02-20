package main

import (
	"context"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/agent"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	cfg := NewConfig()

	if err := logger.Initialize(cfg.flagLogLevel); err != nil {
		fmt.Println("logger initialization failed")
		return
	}

	if err := run(cfg); err != nil {
		logger.Log.Error("While executing agent, error occurred", zap.Error(err))
	}
}

func run(cfg Config) error {

	logger.Log.Info(fmt.Sprintf("Running agent with poll interval %d and report interval %d\n", cfg.pollInterval, cfg.reportInterval))
	logger.Log.Info(fmt.Sprintf("Metric storage server address is set to %s\n", cfg.serverIPAddr))
	a := agent.NewAgent("http://"+cfg.serverIPAddr, 3, 1)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	a.WG.Add(2)
	go a.GetMetrics(ctx, time.Duration(cfg.pollInterval))
	go a.GetGopsutilMetrics(ctx, time.Duration(cfg.pollInterval))

	g := new(errgroup.Group)

	for i := 0; i < int(cfg.flagRateLimit); i++ {
		g.Go(func() error {
			err := a.IterateStructFieldsAndSend(ctx, time.Duration(cfg.reportInterval), cfg.flagSecretKey)
			if err != nil {
				logger.Log.Error("failed to send metrics,", zap.Error(err))
				return err
			}
			return nil
		},
		)
	}
	if err := g.Wait(); err != nil {
		cancel()
	}
	a.WG.Wait()
	return nil
}
