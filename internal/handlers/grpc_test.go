package handlers

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	pb "github.com/sebasttiano/Blackbird.git/internal/proto"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	mockservice "github.com/sebasttiano/Blackbird.git/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
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
					logger.Log.Error("Server exited with error: %v", zap.Error(err))
				}
			}()

			bufDialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}

			ctx := context.TODO()
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				logger.Log.Error("NewClientConn err: %v", zap.Error(err))
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
