package main

import (
	"fmt"
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
		expectedDuration int
	}{
		{
			name:             "Test pollInterval",
			args:             args{pollInterval: 2, reportInterval: 10},
			expectedDuration: 30,
		},
	}

	router := handlers.InitRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	for _, tt := range tests {
		fmt.Println(server)
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			mh := NewMetricHandler(tt.args.pollInterval, tt.args.reportInterval, tt.expectedDuration, "http://"+serverIPAddr)
			mh.GetMetrics()
			duration := time.Since(startTime)
			assert.Equal(t, time.Duration(tt.expectedDuration)*time.Second, duration.Round(time.Second))
		})

	}
}
