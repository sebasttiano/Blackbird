package storage

import (
	"encoding/json"
	"errors"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"github.com/sebasttiano/Blackbird.git/templates"
	"go.uber.org/zap"
	"os"
	"strconv"
)

var SrvFacility = NewServerFacility()

func GetCurrentServerSettings() *ServerSettings {
	return &SrvFacility.Settings
}

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

// GetValue returns either gauge or counter metrics
func (g *MemStorage) GetValue(metricName string, metricType string) (interface{}, error) {
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
func (g *MemStorage) GetModelValue(metric *models.Metrics) error {

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
		return errors.New("error: unknown metric type. Only gauge and counter are available")
	}

	if SrvFacility.Settings.SyncSave {
		if err := SrvFacility.LocalStorage.SaveToFile(SrvFacility.Settings.SaveFilePath); err != nil {
			logger.Log.Error("couldn`t save to the file", zap.Error(err))
			return err
		}
	}
	return nil
}

// SetModelValue saves either gauge or counter metrics from model
func (g *MemStorage) SetModelValue(metric *models.Metrics) error {

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

	if SrvFacility.Settings.SyncSave {
		if err := SrvFacility.LocalStorage.SaveToFile(SrvFacility.Settings.SaveFilePath); err != nil {
			logger.Log.Error("couldn`t save to the file", zap.Error(err))
			return err
		}
	}
	return nil
}

func (g *MemStorage) SaveToFile(file string) error {
	data, err := json.MarshalIndent(g, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0666)
}

func (g *MemStorage) RestoreFromFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, g); err != nil {
		return err
	}
	return nil
}

type HandleMemStorage interface {
	GetValue(metricName string, metricType string) (interface{}, error)
	GetModelValue(metrics *models.Metrics) error
	SetValue(metricName string, metricType string, metricValue string) error
	SetModelValue(metric *models.Metrics) error
	SaveToFile(file string) error
	RestoreFromFile(file string) error
}

type ServerFacility struct {
	LocalStorage  HandleMemStorage
	HTMLTemplates templates.HTMLTemplates
	Settings      ServerSettings
}

func NewServerFacility() ServerFacility {
	return ServerFacility{
		LocalStorage:  NewMemStorage(),
		HTMLTemplates: templates.ParseTemplates(),
		Settings: ServerSettings{
			SyncSave:   false,
			SaveToFile: true,
		},
	}
}

type ServerSettings struct {
	SyncSave     bool
	SaveToFile   bool
	SaveFilePath string
}
