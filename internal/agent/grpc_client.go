package agent

import (
	"context"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	pb "github.com/sebasttiano/Blackbird.git/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GRPCClient реализующий интерфейс Sender, отправляет на gRPC сервер
type GRPCClient struct {
	client pb.MetricsClient
	conn   *grpc.ClientConn
}

// NewGRPCClient - конструктор для GRPCClient
func NewGRPCClient(serverAddr string) (*GRPCClient, error) {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("failed to create grpc client", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrInitSender, err)
	}
	logger.Log.Info("successfully init grpc client", zap.String("address", serverAddr))
	c := pb.NewMetricsClient(conn)

	return &GRPCClient{
		client: c,
		conn:   conn,
	}, nil
}

func (g *GRPCClient) CloseConnection() error {
	err := g.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// SendToRepo собирает из каналов метрики, формирует и шлет protobuf сообщение в репозиторий
func (g *GRPCClient) SendToRepo(jobsMetrics <-chan MetricsSet, jobsGMetrics <-chan GopsutilMetricsSet) error {

	resp, err := g.client.ListAllMetrics(context.Background(), &emptypb.Empty{})
	if err != nil {
		logger.Log.Error("failed to send to repo", zap.Error(err))
	}
	fmt.Println(resp)
	g.CloseConnection()
	return nil
}
