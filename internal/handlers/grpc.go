package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/sebasttiano/Blackbird.git/internal/proto"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ErrInternalGrpc = errors.New("internal grpc server error")

type MetricsServer struct {
	Service *service.Service
	pb.UnimplementedMetricsServer
}

func (m *MetricsServer) ListAllMetrics(ctx context.Context, in *emptypb.Empty) (*pb.ListMetricsResponse, error) {
	var counters []*pb.CounterMetric
	var gauges []*pb.GaugeMetric

	data := m.Service.GetAllValues(ctx)

	jsonGauge, err := json.Marshal(data.Gauge)
	if err != nil {
		return &pb.ListMetricsResponse{}, err
	}

	if err := json.Unmarshal(jsonGauge, &gauges); err != nil {
		return &pb.ListMetricsResponse{}, err
	}

	jsonCounter, err := json.Marshal(data.Counter)
	if err != nil {
		return &pb.ListMetricsResponse{Error: ErrInternalGrpc.Error()}, fmt.Errorf("%w: %v", ErrInternalGrpc, err)
	}

	if err := json.Unmarshal(jsonCounter, &counters); err != nil {
		return &pb.ListMetricsResponse{Error: ErrInternalGrpc.Error()}, fmt.Errorf("%w: %v", ErrInternalGrpc, err)
	}

	response := pb.ListMetricsResponse{
		Counters: counters,
		Gauges:   gauges,
		Error:    "",
	}
	return &response, nil
}
