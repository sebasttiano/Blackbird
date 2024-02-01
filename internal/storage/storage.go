package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"go.uber.org/zap"
	"os"
	"strconv"
)

// MemStorage Keeps Gauge and Counter metrics
type MemStorage struct {
	Gauge    map[string]float64
	Counter  map[string]int64
	Settings *StoreSettings
}

// NewMemStorage — constructor of the type MemStorage.
func NewMemStorage(storeSettings *StoreSettings) *MemStorage {
	return &MemStorage{
		Gauge:    make(map[string]float64),
		Counter:  make(map[string]int64),
		Settings: storeSettings,
	}
}

// GetValue returns either gauge or counter metrics
func (g *MemStorage) GetValue(ctx context.Context, metricName string, metricType string) (interface{}, error) {
	switch metricType {
	case "gauge":
		value, ok := g.Gauge[metricName]
		if !ok {
			return nil, errors.New("error: invalid gauge metric name")
		}
		return value, nil
	case "counter":
		value, ok := g.Counter[metricName]
		if !ok {
			return nil, errors.New("error: invalid counter metric name")
		}
		return value, nil
	default:
		return nil, errors.New("error: unknown metric type. only gauge and counter are available")
	}
}

// GetModelValue returns either gauge or counter metrics
func (g *MemStorage) GetModelValue(ctx context.Context, metric *models.Metrics) error {

	if metric.ID == "" {
		return errors.New("name of the metric is required")
	}

	switch metric.MType {
	case "gauge":

		value, ok := g.Gauge[metric.ID]
		if !ok {
			return errors.New("error: invalid gauge metric name")
		}
		metric.Value = &value
	case "counter":
		sum, ok := g.Counter[metric.ID]
		if !ok {
			return errors.New("error: invalid counter metric name")
		}
		metric.Delta = &sum
	default:
		return errors.New("error: unknown metric type. only gauge and counter are available")
	}
	return nil
}

// SetValue saves either gauge or counter metrics
func (g *MemStorage) SetValue(ctx context.Context, metricName string, metricType string, metricValue string) error {
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
		return errors.New("error: unknown metric type. Only gauge and counter are available")
	}

	if g.Settings.SyncSave {
		if err := g.Save(); err != nil {
			logger.Log.Error("couldn`t save to the file", zap.Error(err))
			return err
		}
	}
	return nil
}

// SetModelValue saves either gauge or counter metrics from model
func (g *MemStorage) SetModelValue(ctx context.Context, metrics []*models.Metrics) error {

	for _, metric := range metrics {
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
			*metric.Delta = g.Counter[metric.ID]
		default:
			return errors.New("error: unknown metric type. Only gauge and counter are available")
		}
	}
	if g.Settings.SyncSave {
		if err := g.Save(); err != nil {
			logger.Log.Error("couldn`t save to the file", zap.Error(err))
			return err
		}
	}
	return nil
}

func (g *MemStorage) Save() error {
	if g.Settings.SaveFilePath == "" {
		return errors.New("can`t save to file. no file path specify")
	}
	data, err := json.MarshalIndent(g, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(g.Settings.SaveFilePath, data, 0666)
}

func (g *MemStorage) Restore() error {
	if g.Settings.SaveFilePath == "" {
		return errors.New("can`t restore from file. no file path specify")
	}
	data, err := os.ReadFile(g.Settings.SaveFilePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, g); err != nil {
		return err
	}
	return nil
}

func (g *MemStorage) GetAllValues(ctx context.Context) *StoreMetrics {

	storeMetrics := &StoreMetrics{}

	for key, value := range g.Gauge {
		storeMetrics.Gauge = append(storeMetrics.Gauge, GaugeMetric{name: key, value: value})
	}

	for key, value := range g.Counter {
		storeMetrics.Counter = append(storeMetrics.Counter, CounterMetric{name: key, value: value})
	}
	return storeMetrics
}
