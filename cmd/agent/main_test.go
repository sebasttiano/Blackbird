package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMetrics(t *testing.T) {
	type args struct {
		pollInterval   int64
		reportInterval int64
	}
	tests := []struct {
		name             string
		args             args
		expectedDuration time.Duration
	}{
		{
			name:             "Test pollInterval",
			args:             args{pollInterval: 2, reportInterval: 10},
			expectedDuration: 30 * time.Second,
		},
	}
	for _, tt := range tests {

		router := handlers.InitRouter()
		server := httptest.NewServer(router)

		httpClient := NewHTTPClient(server.URL)
		defer server.Close()

		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			GetMetrics(tt.args.pollInterval, tt.args.reportInterval, int(tt.expectedDuration.Seconds()), httpClient)
			duration := time.Since(startTime)
			assert.Equal(t, tt.expectedDuration, duration.Round(time.Second))
		})
	}
}
