package storage

import (
	"errors"
	"strconv"
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
func (g *MemStorage) SetValue(metricName string, metricType string, metricValue string) error {
	switch metricType {
	case "gauge":
		valueFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return err
		}
		g.Gauge[metricName] = valueFloat
	case "counter":
		valueInt, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return err
		}
		g.Counter[metricName] += valueInt
	default:
		return errors.New("error: UNKNOWN METRIC TYPE. Only gauge and counter are available")
	}
	return nil
}

type HandleMemStorage interface {
	GetValue(metricName string, metricType string) (interface{}, error)
	SetValue(metricName string, metricType string, metricValue string) error
}
