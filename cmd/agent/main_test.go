package main

import (
	"context"
	"github.com/sebasttiano/Blackbird.git/internal/agent"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMetrics(t *testing.T) {

	views := handlers.NewServerViews(storage.NewMemStorage(&storage.StoreSettings{}))
	router := views.InitRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	a := agent.NewAgent(serverURL, 3, 1)

	t.Run("Test running intervals", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		startTime := time.Now()
		a.WG.Add(1)
		a.GetMetrics(ctx, 1)
		duration := time.Since(startTime)
		assert.Equal(t, time.Duration(5)*time.Second, duration.Round(time.Second))
	})
}

func TestIterateStructFieldsAndSend(t *testing.T) {

	views := handlers.NewServerViews(storage.NewMemStorage(&storage.StoreSettings{}))
	router := views.InitRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	a := agent.NewAgent(serverURL, 3, 1)

	tests := []struct {
		name           string
		notExpectedMsg string
		testMetric     agent.MetricsSet
	}{
		{
			name:           "Test OK return code for all metrics",
			notExpectedMsg: "servers return code",
			testMetric:     agent.MetricsSet{Alloc: 134408, Mallocs: 312, MCacheInuse: 9600},
		},
		{
			name:           "Test nice server parsing",
			notExpectedMsg: "invalid syntax",
			testMetric:     agent.MetricsSet{HeapIdle: 3.35872, NumForcedGC: 0, BuckHashSys: 9600},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Metrics = tt.testMetric
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			if err := a.IterateStructFieldsAndSend(ctx, 5, "secret"); err != nil {
				assert.NotContainsf(t, err.Error(), tt.notExpectedMsg, "not expected error occured")
			}

		})
	}
}
