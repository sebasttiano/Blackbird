package server

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	pb "github.com/sebasttiano/Blackbird.git/internal/proto"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"sync"
)

// GRPSServer реалиузет gRPC сервер.
type GRPSServer struct {
	srv *grpc.Server
}

// NewGRPSServer конструктор для gRPC сервера
func NewGRPSServer(service *service.Service) *GRPSServer {
	s := grpc.NewServer(grpc.UnaryInterceptor(logging.UnaryServerInterceptor(handlers.InterceptorLogger(logger.Log))))
	pb.RegisterMetricsServer(s, &handlers.MetricsServer{Service: service})
	return &GRPSServer{
		srv: s,
	}
}

// Start запускает grpc сервер.
func (s *GRPSServer) Start(addr string) {
	logger.Log.Info("Running gRPC server", zap.String("address", addr))
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Log.Error("failed to allocate tcp socket for gRPC server", zap.Error(err))
	}
	if err := s.srv.Serve(listen); err != nil {
		logger.Log.Error("failed to start gRPC server", zap.Error(err))
	}
}

// HandleShutdown закрывает http сервер.
func (s *GRPSServer) HandleShutdown(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	<-ctx.Done()
	logger.Log.Info("shutdown signal caught. shutting down gRPC server")

	s.srv.GracefulStop()
	logger.Log.Info("gRPC server gracefully shutdown")
}
