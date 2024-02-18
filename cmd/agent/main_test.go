package main

import (
	"context"
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
	mh := NewMetricHandler(1, 5, serverURL, "secret")

	t.Run("Test running intervals", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		startTime := time.Now()
		mh.wg.Add(1)
		mh.GetMetrics(ctx)
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
	mh := NewMetricHandler(1, 5, serverURL, "secret")

	tests := []struct {
		name           string
		notExpectedMsg string
		testMetric     MetricsSet
	}{
		{
			name:           "Test OK return code for all metrics",
			notExpectedMsg: "servers return code",
			testMetric:     MetricsSet{Alloc: 134408, Mallocs: 312, MCacheInuse: 9600},
		},
		{
			name:           "Test nice server parsing",
			notExpectedMsg: "invalid syntax",
			testMetric:     MetricsSet{HeapIdle: 3.35872, NumForcedGC: 0, BuckHashSys: 9600},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh.metrics = tt.testMetric
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			if err := mh.IterateStructFieldsAndSend(ctx); err != nil {
				assert.NotContainsf(t, err.Error(), tt.notExpectedMsg, "not expected error occured")
			}

		})
	}
}
