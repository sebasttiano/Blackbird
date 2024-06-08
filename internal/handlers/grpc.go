package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	pb "github.com/sebasttiano/Blackbird.git/internal/proto"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ErrInternalGrpc = errors.New("internal grpc server error")
var ErrBadRequestGrpc = errors.New("bad request")

type MetricsServer struct {
	Service *service.Service
	pb.UnimplementedMetricsServer
}

// ListAllMetrics возвращает все сохраненные метрики
func (m *MetricsServer) ListAllMetrics(ctx context.Context, in *emptypb.Empty) (*pb.ListMetricsResponse, error) {

	data := m.Service.GetAllValues(ctx)
	var metrics []*pb.Metric

	for _, value := range data.Gauge {
		metrics = append(metrics, &pb.Metric{Id: value.Name, Value: value.Value, Type: pb.MetricType_gauge})
	}

	for _, value := range data.Counter {
		metrics = append(metrics, &pb.Metric{Id: value.Name, Delta: value.Value, Type: pb.MetricType_counter})
	}

	response := pb.ListMetricsResponse{
		Metrics: metrics,
		Error:   "",
	}
	return &response, nil
}

// GetMetric возвращает метрику по переданному типу и имени
func (m *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse

	value, err := m.Service.GetValue(ctx, in.Metric.Id, in.Metric.Type.String())
	if err != nil {
		logger.Log.Error("couldn`t find requested metric. ", zap.Error(err))
		response.Error = fmt.Sprintf("couldn`t find requested metric, %s", in.Metric.Id)
		return &response, nil
	}

	response.Metric = in.Metric

	switch value := value.(type) {
	case float64:
		response.Metric.Value = value
		response.Metric.Type = pb.MetricType_gauge
	case int64:
		response.Metric.Delta = value
		response.Metric.Type = pb.MetricType_counter
	}

	return &response, nil

}

// UpdateMetric обновляет одну метрику
func (m *MetricsServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse

	if err := m.Service.SetValue(ctx, in.Id, in.Type.String(), in.Value); err != nil {
		logger.Log.Error("couldn`t save metric. error: ", zap.Error(err))
		response.Error = fmt.Sprintf("error: couldn`t save metric, %s", in.Id)
	}
	return &response, nil
}

// UpdateMetrics обновляет сет метрик
func (m *MetricsServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricResponse, error) {
	var metricSet models.MetricSet

	marshaller := protojson.MarshalOptions{EmitDefaultValues: true}
	jsonMetrics, err := marshaller.Marshal(in)
	if err != nil {
		logger.Log.Error("failed to marshal metrics to json", zap.Error(err))
		return &pb.UpdateMetricResponse{Error: ErrInternalGrpc.Error()}, fmt.Errorf("%w: %v", ErrInternalGrpc, err)
	}

	if err := json.Unmarshal(jsonMetrics, &metricSet); err != nil {
		logger.Log.Error("couldn`t unmarshal json metrics", zap.Error(err))
		return &pb.UpdateMetricResponse{Error: ErrInternalGrpc.Error()}, fmt.Errorf("%w: %v", ErrInternalGrpc, err)
	}

	if err := m.Service.SetModelValue(ctx, metricSet.CastToMetrics()); err != nil {
		logger.Log.Error("couldn`t save metric. error: ", zap.Error(err))
		return &pb.UpdateMetricResponse{Error: ErrInternalGrpc.Error()}, fmt.Errorf("%w: %v", ErrInternalGrpc, err)
	}
	return &pb.UpdateMetricResponse{}, nil
}
