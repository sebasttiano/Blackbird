package agent

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetMetrics(t *testing.T) {
	views := handlers.NewServerViews(service.NewService(&service.Settings{}, repository.NewMemStorage()))
	router := views.InitRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	a, _ := NewAgent(serverURL, 3, 1, "", nil, "")

	t.Run("Test running intervals", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		startTime := time.Now()
		a.WG.Add(1)
		a.GetMetrics(ctx, 1, make(chan<- MetricsSet))
		duration := time.Since(startTime)
		assert.Equal(t, time.Duration(5)*time.Second, duration.Round(time.Second))
	})
}

func BenchmarkAgentMetrics(b *testing.B) {
	a, _ := NewAgent("localhost:8080", 1, 1, "", nil, "")

	var jobsMetricCount int
	var jobsGMetricCount int

	b.Run("BenchmarkGetMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(100*b.N)*time.Millisecond)
			jobsMetric := make(chan MetricsSet, 50)
			a.WG.Add(1)
			a.GetMetrics(ctx, 10*time.Millisecond, jobsMetric)
			jobsMetricCount = +len(jobsMetric)
			cancel()
		}
		b.ReportMetric(float64(jobsMetricCount)/float64(b.N), "metric_jobs/op")
	})
	b.Run("BenchmarkGetGopsutilMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(100*b.N)*time.Millisecond)
			jobsGMetrics := make(chan GopsutilMetricsSet, 100)
			a.WG.Add(1)
			a.GetGopsutilMetrics(ctx, 10*time.Millisecond, jobsGMetrics)
			jobsGMetricCount += len(jobsGMetrics)
			cancel()
		}
		b.ReportMetric(float64(jobsGMetricCount)/float64(b.N), "gopsutil_jobs/op")
	})
}
