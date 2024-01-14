package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
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

	ticker := time.NewTicker(time.Second * 5)
	go common.Schedule(ticker)

	return http.ListenAndServe(flagRunAddr, handlers.WithLogging(handlers.GzipMiddleware(handlers.InitRouter())))
}
