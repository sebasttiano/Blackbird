package repository

import (
	"context"
	"errors"
)

// MemStorage хранит Gauge и Counter метрики в памяти
type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// NewMemStorage конструктор для MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

// GetGauge метод из памяти возвращает сохраненную метрику типа Gauge.
func (g *MemStorage) GetGauge(ctx context.Context, metric *GaugeMetric) error {
	var ok bool
	metric.Value, ok = g.Gauge[metric.Name]
	if !ok {
		return errors.New("error: invalid gauge metric name")
	}
	return nil
}

// GetCounter метод из памяти возвращает сохраненную метрику типа Counter
func (g *MemStorage) GetCounter(ctx context.Context, metric *CounterMetric) error {
	var ok bool
	metric.Value, ok = g.Counter[metric.Name]
	if !ok {
		return errors.New("error: invalid counter metric name")
	}
	return nil
}

// SetGauge метод сохраняет в памяти метрику типа Gauge.
func (g *MemStorage) SetGauge(ctx context.Context, metric *GaugeMetric) error {
	g.Gauge[metric.Name] = metric.Value
	return nil
}

// SetCounter метод сохоаняет в памяти метрику типа Counter.
func (g *MemStorage) SetCounter(ctx context.Context, metric *CounterMetric) error {
	g.Counter[metric.Name] += metric.Value
	return nil
}

// GetAllMetrics метод возвращает все метрики из памяти.
func (g *MemStorage) GetAllMetrics(ctx context.Context, s *StoreMetrics) error {
	for key, value := range g.Gauge {
		s.Gauge = append(s.Gauge, GaugeMetric{Name: key, Value: value})
	}

	for key, value := range g.Counter {
		s.Counter = append(s.Counter, CounterMetric{Name: key, Value: value})
	}
	return nil
}

// RestoreAllMetrics восстанавливает в памяти все метрики.
func (g *MemStorage) RestoreAllMetrics(gauges map[string]float64, counters map[string]int64) {
	g.Gauge = gauges
	g.Counter = counters
}
