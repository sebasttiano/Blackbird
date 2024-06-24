package server

import (
	"context"
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/sebasttiano/Blackbird.git/internal/config"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/logger"

	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Server http
type Server struct {
	srv   *http.Server
	views *handlers.ServerViews
}

// NewServer конструктор для типа http сервер.
func NewServer(serverAddr string, views *handlers.ServerViews, router chi.Router) *Server {

	return &Server{
		srv:   &http.Server{Addr: serverAddr, Handler: router},
		views: views,
	}
}

// Start запускает http сервер.
func (s *Server) Start(cfg *config.Config) {
	logger.Log.Info("Running server", zap.String("address", cfg.ServerIPAddr))
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Error("server error", zap.Error(err))
		return
	}
}

// HandleShutdown закрывает http сервер.
func (s *Server) HandleShutdown(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) {

	defer wg.Done()

	<-ctx.Done()
	logger.Log.Info("shutdown signal caught. shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if cfg.FileStoragePath != "" {
		if err := s.views.Service.Save(); err != nil {
			logger.Log.Error("couldn`t finally save file after graceful shutdown", zap.Error(err))
		}
	}
	err := s.srv.Shutdown(ctx)
	if err != nil {
		logger.Log.Error("server shutdown error", zap.Error(err))
		return
	}
	logger.Log.Info("server gracefully shutdown")
}
