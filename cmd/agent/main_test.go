package main

import (
	"context"
	"github.com/sebasttiano/Blackbird.git/internal/agent"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetMetrics(t *testing.T) {

	views := handlers.NewServerViews(repository.NewMemStorage(&repository.StoreSettings{}))
	router := views.InitRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	a := agent.NewAgent(serverURL, 3, 1, "")

	t.Run("Test running intervals", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		startTime := time.Now()
		a.WG.Add(1)
		a.GetMetrics(ctx, 1, make(chan<- agent.MetricsSet))
		duration := time.Since(startTime)
		assert.Equal(t, time.Duration(5)*time.Second, duration.Round(time.Second))
	})
}
