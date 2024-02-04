package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// DBStorage Keeps metrics in database
type DBStorage struct {
	conn   *sqlx.DB
	client *PGClient
}

type DBErrors struct {
	ErrConnect error
}

var PGError *pgconn.PgError

// GetValue returns either gauge or counter metrics
func (d *DBStorage) GetValue(ctx context.Context, metricName string, metricType string) (interface{}, error) {

	switch metricType {
	case "gauge":
		m := GaugeMetric{Name: metricName}
		g, err := Retry(ctx, d.client.retryDelays, d.client.GetGauge, &m)
		if err != nil {
			return nil, err
		}
		return g.Value, nil
	case "counter":
		m := CounterMetric{Name: metricName}
		c, err := Retry(ctx, d.client.retryDelays, d.client.GetCounter, &m)
		if err != nil {
			return nil, err
		}
		return c.Value, nil
	default:
		return nil, errors.New("error: unknown metric type. only gauge and counter are available")
	}
}

// GetModelValue returns either gauge or counter metrics
func (d *DBStorage) GetModelValue(ctx context.Context, metric *models.Metrics) error {

	if metric.ID == "" {
		return errors.New("name of the metric is required")
	}

	value, err := d.GetValue(ctx, metric.ID, metric.MType)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case float64:
		metric.Value = &v
	case int64:
		metric.Delta = &v
	default:
		return errors.New("error: unknown metric type. only gauge and counter are available")
	}
	return nil
}

// SetValue saves either gauge or counter metrics
func (d *DBStorage) SetValue(ctx context.Context, metricName string, metricType string, metricValue string) error {

	switch metricType {
	case "gauge":
		valueFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return err
		}
		m := GaugeMetric{Name: metricName, Value: valueFloat}
		_, err = Retry(ctx, d.client.retryDelays, d.client.SetGauge, &m)
		if err != nil {
			return err
		}
	case "counter":
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return err
		}
		m := CounterMetric{Name: metricName, Value: intValue}
		_, err = Retry(ctx, d.client.retryDelays, d.client.SetCounter, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetModelValue saves either gauge or counter metrics from model
func (d *DBStorage) SetModelValue(ctx context.Context, metrics []*models.Metrics) error {

	for _, metric := range metrics {

		if metric.ID == "" {
			return errors.New("name of the metric is required")
		}

		switch metric.MType {
		case "gauge":
			if metric.Value == nil {
				return errors.New("value of the gauge is required")
			}
			if err := d.SetValue(ctx, metric.ID, metric.MType, fmt.Sprintf("%.12f", *metric.Value)); err != nil {
				return err
			}

		case "counter":
			if metric.Delta == nil {
				return errors.New("value of the counter is required")
			}
			if err := d.SetValue(ctx, metric.ID, metric.MType, fmt.Sprintf("%d", *metric.Delta)); err != nil {
				return err
			}
		default:
			return errors.New("error: unknown metric type. Only gauge and counter are available")
		}
	}
	return nil
}

// GetAllValues get all metrics from db and returns in raw format
func (d *DBStorage) GetAllValues(ctx context.Context) (s *StoreMetrics) {

	s = &StoreMetrics{make([]GaugeMetric, 0), make([]CounterMetric, 0)}

	Retry(ctx, d.client.retryDelays, d.client.GetAllMetrics, s)
	return s
}

// NewDBStorage returns new database storage
func NewDBStorage(conn *sqlx.DB, bootstrap bool, retries uint, backoffFactor uint) *DBStorage {
	db := &DBStorage{conn: conn}
	if bootstrap {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := db.Bootstrap(ctx); err != nil {
			if errors.As(err, &PGError) {
				if PGError.Code == pgerrcode.InFailedSQLTransaction {
					logger.Log.Debug("rollback in bootstrap occured!")
				} else {
					logger.Log.Error("db bootstrap failed", zap.Error(err))
				}
			}
		}
	}
	return &DBStorage{conn: conn, client: NewPGClient(conn, retries, backoffFactor)}
}

func (d *DBStorage) Save() error {
	return nil
}

func (d *DBStorage) Restore() error {
	return nil
}

// Bootstrap creates tables in DB
func (d *DBStorage) Bootstrap(ctx context.Context) error {

	logger.Log.Debug("checking db tables")
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// create table for gauge metrics
	if _, err := tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS gauge_metrics (
            id serial PRIMARY KEY,
			name varchar(128),
            gauge double precision 
        )
	`); err != nil {
		logger.Log.Error("failed to create gauge_metrics table", zap.Error(err))
		return err
	}

	// create table for counter metrics
	if _, err := tx.ExecContext(ctx, `
	   CREATE TABLE IF NOT EXISTS counter_metrics (
	       id serial PRIMARY KEY,
		   name varchar(128),
	       counter bigint
	   )
	`); err != nil {
		logger.Log.Error("failed to create gauge_metrics table", zap.Error(err))
		return err
	}

	// commit
	return tx.Commit()
}
