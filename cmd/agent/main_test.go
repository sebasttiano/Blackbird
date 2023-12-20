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
		name string
		args args
	}{
		{
			name: "Test pollInterval",
			args: args{pollInterval: 2, reportInterval: 10},
		},
	}
	for _, tt := range tests {

		router := handlers.InitRouter()
		server := httptest.NewServer(router)

		httpClient := NewHTTPClient(server.URL)
		defer server.Close()

		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now().Unix()
			GetMetrics(tt.args.pollInterval, tt.args.reportInterval, int(tt.args.reportInterval/tt.args.pollInterval), httpClient)
			endTime := time.Now().Unix()
			assert.Equal(t, reportInterval, endTime-startTime)
		})
	}
}
