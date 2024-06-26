package handlers

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	pb "github.com/sebasttiano/Blackbird.git/internal/proto"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	mockservice "github.com/sebasttiano/Blackbird.git/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func TestMetricsServer_GetMetric(t *testing.T) {

	type mockBehaviour func(s *mockservice.MockMetricService, in *pb.GetMetricRequest)

	testTable := []struct {
		name          string
		in            *pb.GetMetricRequest
		mockBehaviour mockBehaviour
		expected      *pb.GetMetricResponse
		err           error
	}{
		{
			name: "OK counter metric",
			in: &pb.GetMetricRequest{
				Metric: &pb.Metric{Id: "test_counter", Delta: 0, Value: 0, Type: pb.MetricType_counter},
			},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.GetMetricRequest) {
				var value int64 = 33
				s.EXPECT().GetValue(gomock.Any(), in.Metric.Id, in.Metric.Type.String()).Return(value, nil)
			},
			expected: &pb.GetMetricResponse{
				Metric: &pb.Metric{Id: "test_counter", Delta: 33},
			},
			err: nil,
		},
		{
			name: "NOT OK, unknown metric type",
			in: &pb.GetMetricRequest{
				Metric: &pb.Metric{Id: "test_counter", Delta: 0, Value: 0, Type: pb.MetricType_counter},
			},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.GetMetricRequest) {
				s.EXPECT().GetValue(gomock.Any(), in.Metric.Id, in.Metric.Type.String()).Return(nil, service.ErrUnknownMetricType)
			},
			expected: nil,
			err:      status.Errorf(codes.InvalidArgument, "invalid argument: test_counter - counter"),
		},
		{
			name: "NOT OK, metric not found",
			in: &pb.GetMetricRequest{
				Metric: &pb.Metric{Id: "alloc", Delta: 0, Value: 0, Type: pb.MetricType_counter},
			},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.GetMetricRequest) {
				s.EXPECT().GetValue(gomock.Any(), in.Metric.Id, in.Metric.Type.String()).Return(nil, errors.New("metric not found"))
			},
			expected: nil,
			err:      status.Errorf(codes.NotFound, "couldn`t find requested metric. alloc"),
		},
		{
			name: "OK gauge metric",
			in: &pb.GetMetricRequest{
				Metric: &pb.Metric{Id: "test_gauge", Delta: 0, Value: 0, Type: pb.MetricType_gauge},
			},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.GetMetricRequest) {
				value := 197.30
				s.EXPECT().GetValue(gomock.Any(), in.Metric.Id, in.Metric.Type.String()).Return(value, nil)
			},
			expected: &pb.GetMetricResponse{
				Metric: &pb.Metric{Id: "test_gauge", Value: 197.30},
			},
			err: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()

			mock := mockservice.NewMockMetricService(c)
			tt.mockBehaviour(mock, tt.in)

			pb.RegisterMetricsServer(s, &MetricsServer{Service: mock})
			go func() {
				if err := s.Serve(lis); err != nil {
					t.Errorf("Server exited with error: %v", err)
				}
			}()

			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}

			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("NewClientConn err: %v", err)
			}
			defer conn.Close()
			client := pb.NewMetricsClient(conn)

			resp, err := client.GetMetric(ctx, tt.in)
			if tt.err != nil {
				assert.Errorf(t, err, tt.err.Error())
				assert.Equal(t, tt.err, err)
			} else {
				assert.Equal(t, resp.Metric.Id, tt.expected.Metric.Id)
				assert.Equal(t, resp.Metric.Delta, tt.expected.Metric.Delta)
				assert.Equal(t, resp.Metric.Value, tt.expected.Metric.Value)
			}

		})
	}
}

func TestMetricsServer_UpdateMetric(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockMetricService, in *pb.UpdateMetricRequest)

	testTable := []struct {
		name          string
		in            *pb.UpdateMetricRequest
		mockBehaviour mockBehaviour
		err           error
	}{
		{
			name: "Ok counter metric",
			in:   &pb.UpdateMetricRequest{Id: "test_counter", Value: "100", Type: pb.MetricType_counter},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.UpdateMetricRequest) {
				s.EXPECT().SetValue(gomock.Any(), in.Id, in.Type.String(), in.Value).Return(nil)
			},
			err: nil,
		},
		{
			name: "Ok gauge metric",
			in:   &pb.UpdateMetricRequest{Id: "test_gauge", Value: "33.313", Type: pb.MetricType_gauge},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.UpdateMetricRequest) {
				s.EXPECT().SetValue(gomock.Any(), in.Id, in.Type.String(), in.Value).Return(nil)
			},
			err: nil,
		},
		{
			name: "NOT OK, unknown metric type",
			in:   &pb.UpdateMetricRequest{Id: "test_gauge", Value: "33.313", Type: pb.MetricType_gauge},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.UpdateMetricRequest) {
				s.EXPECT().SetValue(gomock.Any(), in.Id, in.Type.String(), in.Value).Return(service.ErrUnknownMetricType)
			},
			err: status.Errorf(codes.InvalidArgument, "invalid argument: test_gauge - gauge"),
		},
		{
			name: "NOT OK, failed to save metric",
			in:   &pb.UpdateMetricRequest{Id: "test_gauge", Value: "33.313", Type: pb.MetricType_gauge},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.UpdateMetricRequest) {
				s.EXPECT().SetValue(gomock.Any(), in.Id, in.Type.String(), in.Value).Return(errors.New("failed to save metric"))
			},
			err: status.Errorf(codes.Unknown, "failed to save metric: test_gauge"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()

			mock := mockservice.NewMockMetricService(c)
			tt.mockBehaviour(mock, tt.in)

			pb.RegisterMetricsServer(s, &MetricsServer{Service: mock})
			go func() {
				if err := s.Serve(lis); err != nil {
					t.Errorf("Server exited with error: %v", err)
				}
			}()

			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}

			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("NewClientConn err: %v", err)
			}
			defer conn.Close()
			client := pb.NewMetricsClient(conn)

			_, err = client.UpdateMetric(ctx, tt.in)
			if tt.err != nil {
				assert.Errorf(t, err, tt.err.Error())
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetricsServer_UpdateMetrics(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockMetricService, in *pb.UpdateMetricsRequest)

	testTable := []struct {
		name          string
		in            *pb.UpdateMetricsRequest
		mockBehaviour mockBehaviour
		err           error
	}{
		{
			name: "Ok update metrics",
			in: &pb.UpdateMetricsRequest{Metrics: []*pb.Metric{
				{Id: "test_gauge", Delta: 0, Value: -100.33, Type: pb.MetricType_gauge},
				{Id: "test_counter", Delta: 30, Value: 0, Type: pb.MetricType_counter}}},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.UpdateMetricsRequest) {
				s.EXPECT().SetModelValue(gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "NOT OK, unknown metric type",
			in: &pb.UpdateMetricsRequest{Metrics: []*pb.Metric{
				{Id: "test_gauge", Delta: 0, Value: 0, Type: pb.MetricType_gauge},
				{Id: "test_counter", Delta: 30, Value: 0, Type: pb.MetricType_counter}}},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.UpdateMetricsRequest) {
				s.EXPECT().SetModelValue(gomock.Any(), gomock.Any()).Return(service.ErrUnknownMetricType)
			},
			err: status.Errorf(codes.InvalidArgument, "invalid argument"),
		},
		{
			name: "NOT OK, unknown metric type",
			in: &pb.UpdateMetricsRequest{Metrics: []*pb.Metric{
				{Id: "test_gauge", Delta: 0, Value: 0, Type: pb.MetricType_gauge},
				{Id: "test_counter", Delta: 30, Value: 0, Type: pb.MetricType_counter}}},
			mockBehaviour: func(s *mockservice.MockMetricService, in *pb.UpdateMetricsRequest) {
				s.EXPECT().SetModelValue(gomock.Any(), gomock.Any()).Return(errors.New("failed to save metric"))
			},
			err: status.Errorf(codes.Unknown, "failed to save metrics"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()

			mock := mockservice.NewMockMetricService(c)
			tt.mockBehaviour(mock, tt.in)

			pb.RegisterMetricsServer(s, &MetricsServer{Service: mock})
			go func() {
				if err := s.Serve(lis); err != nil {
					t.Errorf("Server exited with error: %v", err)
				}
			}()

			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}

			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("NewClientConn err: %v", err)
			}
			defer conn.Close()
			client := pb.NewMetricsClient(conn)

			_, err = client.UpdateMetrics(ctx, tt.in)
			if tt.err != nil {
				assert.Errorf(t, err, tt.err.Error())
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetricsServer_ListAllMetrics(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockMetricService)

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		expected      *pb.ListMetricsResponse
	}{
		{
			name: "OK list metrics",
			mockBehaviour: func(s *mockservice.MockMetricService) {
				s.EXPECT().GetAllValues(gomock.Any()).Return(&repository.StoreMetrics{
					Gauge:   []repository.GaugeMetric{{ID: 3, Name: "test_gauge", Value: 29.87}},
					Counter: []repository.CounterMetric{{ID: 19, Name: "test_counter", Value: 99}},
				})
			},
			expected: &pb.ListMetricsResponse{Metrics: []*pb.Metric{
				{Id: "test_gauge", Delta: 0, Value: 29.87, Type: pb.MetricType_gauge},
				{Id: "test_counter", Delta: 30, Value: 0, Type: pb.MetricType_counter},
			}},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer()

			mock := mockservice.NewMockMetricService(c)
			tt.mockBehaviour(mock)

			pb.RegisterMetricsServer(s, &MetricsServer{Service: mock})
			go func() {
				if err := s.Serve(lis); err != nil {
					t.Errorf("Server exited with error: %v", err)
				}
			}()

			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}

			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Errorf("NewClientConn err: %v", err)
			}
			defer conn.Close()
			client := pb.NewMetricsClient(conn)

			var emptyReq emptypb.Empty
			resp, err := client.ListAllMetrics(ctx, &emptyReq)
			assert.NoError(t, err)
			for i, metric := range resp.Metrics {
				assert.Equal(t, metric.Id, tt.expected.Metrics[i].Id)
				assert.Equal(t, metric.Value, tt.expected.Metrics[i].Value)
				assert.Equal(t, metric.Type, tt.expected.Metrics[i].Type)
			}
		})
	}
}
