package repository

import (
	"context"
	"errors"
)

// MemStorage Keeps Gauge and Counter metrics
type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// NewMemStorage — constructor of the type MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

// GetGauge returns gauge  metrics
func (g *MemStorage) GetGauge(ctx context.Context, metric *GaugeMetric) error {
	var ok bool
	metric.Value, ok = g.Gauge[metric.Name]
	if !ok {
		return errors.New("error: invalid gauge metric name")
	}
	return nil
}

// GetCounter returns gauge  metrics
func (g *MemStorage) GetCounter(ctx context.Context, metric *CounterMetric) error {
	var ok bool
	metric.Value, ok = g.Counter[metric.Name]
	if !ok {
		return errors.New("error: invalid counter metric name")
	}
	return nil
}

func (g *MemStorage) SetGauge(ctx context.Context, metric *GaugeMetric) error {
	g.Gauge[metric.Name] = metric.Value
	return nil
}

func (g *MemStorage) SetCounter(ctx context.Context, metric *CounterMetric) error {
	g.Counter[metric.Name] = metric.Value
	return nil
}

func (g *MemStorage) GetAllMetrics(ctx context.Context, s *StoreMetrics) error {

	for key, value := range g.Gauge {
		s.Gauge = append(s.Gauge, GaugeMetric{Name: key, Value: value})
	}

	for key, value := range g.Counter {
		s.Counter = append(s.Counter, CounterMetric{Name: key, Value: value})
	}
	return nil
}

func (g *MemStorage) RestoreAllMetrics(gauges map[string]float64, counters map[string]int64) {
	g.Gauge = gauges
	g.Counter = counters
}
