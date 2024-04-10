package agent

import (
	"context"
	"testing"
	"time"
)

func BenchmarkAgentMetrics(b *testing.B) {
	a := NewAgent("localhost:8080", 1, 1, "")
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	jobsMetric := make(chan MetricsSet)
	jobsGMetrics := make(chan GopsutilMetricsSet)

	var jobsMetricCount int
	var jobsGMetricCount int

	b.Run("BenchmarkGetMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a.WG.Add(1)
			go a.GetMetrics(ctx, 1*time.Second, jobsMetric)
			for range jobsMetric {
				jobsMetricCount++
			}
		}
		b.ReportMetric(float64(jobsMetricCount)/float64(b.N), "metric_jobs/op")
	})
	b.Run("BenchmarkGetGopsutilMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a.WG.Add(1)
			go a.GetGopsutilMetrics(ctx, 1*time.Second, jobsGMetrics)
			for range jobsGMetrics {
				jobsGMetricCount++
			}
		}
		b.ReportMetric(float64(jobsGMetricCount)/float64(b.N), "gopsutil_jobs/op")
	})
	a.WG.Wait()
}
