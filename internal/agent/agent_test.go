package agent

import (
	"context"
	"testing"
	"time"
)

func BenchmarkAgentMetrics(b *testing.B) {
	a := NewAgent("localhost:8080", 1, 1, "", "")

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
