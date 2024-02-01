package storage

import (
	"context"
	"database/sql"
	"github.com/sebasttiano/Blackbird.git/internal/models"
)

type StoreSettings struct {
	SyncSave     bool
	FileSave     bool
	DBSave       bool
	Conn         *sql.DB
	SaveFilePath string
}

type GaugeMetric struct {
	name  string
	value float64
}

type CounterMetric struct {
	name  string
	value int64
}

type StoreMetrics struct {
	Gauge   []GaugeMetric
	Counter []CounterMetric
}

type Store interface {
	GetValue(ctx context.Context, string, metricType string) (interface{}, error)
	GetModelValue(ctx context.Context, metric *models.Metrics) error
	SetValue(ctx context.Context, metricName string, metricType string, metricValue string) error
	SetModelValue(ctx context.Context, metrics []*models.Metrics) error
	GetAllValues(ctx context.Context) *StoreMetrics
	Save() error
	Restore() error
}
