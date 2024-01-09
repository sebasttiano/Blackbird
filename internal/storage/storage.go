package storage

import (
	"errors"
	"github.com/sebasttiano/Blackbird.git/internal/models"
)

// MemStorage Keeps Gauge and Counter metrics
type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// NewMemStorage — constructor of the type MemStorage.
func NewMemStorage() MemStorage {
	return MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

// GetValue returns either gauge or counter metrics
func (g *MemStorage) GetValue(metricName string, metricType string) (interface{}, error) {
	switch metricType {
	case "gauge":
		value, ok := g.Gauge[metricName]
		if !ok {
			return nil, errors.New("error: invalid Gauge metric name")
		}
		return value, nil
	case "counter":
		value, ok := g.Counter[metricName]
		if !ok {
			return nil, errors.New("error: invalid Counter metric name")
		}
		return value, nil
	default:
		return nil, errors.New("error: unknown metric type. only gauge and counter are available")
	}
}

// SetValue saves either gauge or counter metrics
func (g *MemStorage) SetValue(metric models.Metrics) error {

	if metric.ID == "" {
		return errors.New("name of the metric is required")
	}

	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			return errors.New("value of the gauge is required")
		}

		g.Gauge[metric.ID] = *metric.Value
	case "counter":

		if metric.Delta == nil {
			return errors.New("value of the gauge is required")
		}
		g.Counter[metric.ID] += *metric.Delta
	default:
		return errors.New("error: unknown metric type. Only gauge and counter are available")
	}
	return nil
}

type HandleMemStorage interface {
	GetValue(metricName string, metricType string) (interface{}, error)
	SetValue(metric models.Metrics) error
}
