package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMetrics(t *testing.T) {

	serverIPAddr = "localhost:8080"
	type args struct {
		pollInterval   int64
		reportInterval int64
		stopLimit      int
	}
	tests := []struct {
		name                  string
		args                  args
		expectedReturnedConde int
	}{
		{
			name:                  "Test OK return code for all metrics",
			args:                  args{pollInterval: 1, reportInterval: 2, stopLimit: 2},
			expectedReturnedConde: http.StatusOK,
		},
	}

	router := handlers.InitRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	t.Run("Test running intervals", func(t *testing.T) {
		startTime := time.Now()
		mh := NewMetricHandler(2, 10, 30, serverURL)
		responses := mh.GetMetrics()
		for _, resp := range responses {
			resp.Body.Close()
		}

		duration := time.Since(startTime)
		assert.Equal(t, time.Duration(30)*time.Second, duration.Round(time.Second))
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh := NewMetricHandler(tt.args.pollInterval, tt.args.reportInterval, tt.args.stopLimit, serverURL)
			responses := mh.GetMetrics()
			for _, resp := range responses {
				assert.Equal(t, tt.expectedReturnedConde, resp.StatusCode)
				resp.Body.Close()
			}
		})
	}

}
