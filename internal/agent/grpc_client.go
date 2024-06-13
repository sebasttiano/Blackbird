package agent

import (
	"context"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	pb "github.com/sebasttiano/Blackbird.git/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"reflect"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var metric MetricsSet
	var metricG GopsutilMetricsSet
	var metricsBatch []*pb.Metric
	var value reflect.Value

	select {
	case metric = <-jobsMetrics:
		value = reflect.ValueOf(metric)
	case metricG = <-jobsGMetrics:
		value = reflect.ValueOf(metricG)
	}
	numFields := value.NumField()
	structType := value.Type()

	for i := 0; i < numFields; i++ {
		var metrics pb.Metric
		field := structType.Field(i)
		fieldValue := value.Field(i)
		metrics.Id = field.Name

		if fieldValue.CanInt() {
			counterVal := fieldValue.Int()
			metrics.Delta = counterVal
			metrics.Type = pb.MetricType_counter
		} else {
			gaugeVal := fieldValue.Float()
			metrics.Value = gaugeVal
			metrics.Type = pb.MetricType_gauge
		}
		if metrics.Id == "GCCPUFraction" {
			fmt.Println(metrics.Delta)
		}
		metricsBatch = append(metricsBatch, &metrics)
	}

	if len(metricsBatch) > 0 {
		_, err := g.client.UpdateMetrics(ctx, &pb.UpdateMetricsRequest{Metrics: metricsBatch})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.DeadlineExceeded:
					logger.Log.Error("server context deadline exceeded", zap.String("error", e.Message()))
				case codes.InvalidArgument:
					logger.Log.Error("server couldn`t parse metrics", zap.String("server error:", e.Message()))
				default:
					logger.Log.Error("server error", zap.String("error", e.Message()))
				}
			} else {
				logger.Log.Error("failed to update metrics", zap.Error(err))
			}
			return err
		}
		logger.Log.Info("send metrics to repository server successfully.")
	}

	return nil
}
