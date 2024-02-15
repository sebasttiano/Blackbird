package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMetrics(t *testing.T) {

	views := handlers.NewServerViews(storage.NewMemStorage(&storage.StoreSettings{}))
	router := handlers.GzipMiddleware(views.InitRouter())
	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	t.Run("Test running intervals", func(t *testing.T) {
		startTime := time.Now()
		mh := NewMetricHandler(1, 5, 15, serverURL)
		err := mh.GetMetrics()
		require.NoError(t, err)
		duration := time.Since(startTime)
		assert.Equal(t, time.Duration(15)*time.Second, duration.Round(time.Second))
	})
}

func TestIterateStructFieldsAndSend(t *testing.T) {

	views := handlers.NewServerViews(storage.NewMemStorage(&storage.StoreSettings{}))
	router := handlers.GzipMiddleware(views.InitRouter())
	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	client := common.NewHTTPClient(serverURL, 3, 3)

	tests := []struct {
		name           string
		notExpectedMsg string
		testMetric     MetricsSet
	}{
		{
			name:           "Test OK return code for all metrics",
			notExpectedMsg: "server return code",
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
			if err := IterateStructFieldsAndSend(tt.testMetric, client); err != nil {
				assert.NotContainsf(t, err.Error(), tt.notExpectedMsg, "not expected error occured")
			}

		})
	}
}
