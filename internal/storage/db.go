package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"go.uber.org/zap"
	"time"
)

// DBStorage Keeps metrics in database
type DBStorage struct {
	conn *sql.DB
}

var pgError *pgconn.PgError

// GetValue returns either gauge or counter metrics
func (d *DBStorage) GetValue(ctx context.Context, metricName string, metricType string) (interface{}, error) {

	switch metricType {
	case "gauge":
		var m GaugeMetric
		row := d.conn.QueryRowContext(ctx, `SELECT gauge FROM gauge_metrics WHERE name = $1`, metricName)
		if err := row.Scan(&m.value); err != nil {
			return nil, err
		}
		return m.value, nil
	case "counter":
		var m CounterMetric
		row := d.conn.QueryRowContext(ctx, `SELECT counter FROM counter_metrics WHERE name = $1`, metricName)
		if err := row.Scan(&m.value); err != nil {
			return nil, err
		}
		return m.value, nil
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
		var m GaugeMetric
		row := d.conn.QueryRowContext(ctx, `SELECT name, gauge FROM gauge_metrics WHERE name = $1`, metricName)
		if err := row.Scan(&m.name, &m.value); err != nil {
			if _, err := d.conn.ExecContext(ctx, `
			   INSERT INTO gauge_metrics
			   (name, gauge)
			   VALUES
			   ($1, $2);
			`, metricName, metricValue); err != nil {
				return err
			}
		} else {
			if _, err := d.conn.ExecContext(ctx, `
				UPDATE gauge_metrics
				SET gauge = $1
				WHERE name = $2;
			`, metricValue, m.name); err != nil {
				return err
			}
		}

	case "counter":
		var m CounterMetric
		row := d.conn.QueryRowContext(ctx, `SELECT name, counter FROM counter_metrics WHERE name = $1`, metricName)
		if err := row.Scan(&m.name, &m.value); err != nil {
			if _, err := d.conn.ExecContext(ctx, `
			   INSERT INTO counter_metrics
			   (name, counter)
			   VALUES
			   ($1, $2);
			`, metricName, metricValue); err != nil {
				return err
			}
		} else {
			if _, err := d.conn.ExecContext(ctx, `
				UPDATE counter_metrics
				SET counter = counter + $1
				WHERE name = $2;
			`, metricValue, m.name); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetModelValue saves either gauge or counter metrics from model
func (d *DBStorage) SetModelValue(ctx context.Context, metric *models.Metrics) error {

	if metric.ID == "" {
		return errors.New("name of the metric is required")
	}

	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			return errors.New("value of the gauge is required")
		}
		if err := d.SetValue(ctx, metric.ID, metric.MType, fmt.Sprintf("%f", *metric.Value)); err != nil {
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
	return nil
}

// GetAllValues get all metrics from db and returns in raw format
func (d *DBStorage) GetAllValues(ctx context.Context) (s *StoreMetrics) {

	s = &StoreMetrics{make([]GaugeMetric, 0), make([]CounterMetric, 0)}

	gRows, err := d.conn.QueryContext(ctx, `
	   SELECT
	       g.name,
	       g.gauge
	   FROM gauge_metrics g`,
	)

	if err != nil {
		logger.Log.Error("failed to download gauge metrics")
	} else {
		for gRows.Next() {
			var m GaugeMetric
			if err := gRows.Scan(&m.name, &m.value); err != nil {
				logger.Log.Error("while reading from db error occured:", zap.Error(err))
			}
			s.Gauge = append(s.Gauge, m)
		}
	}
	defer gRows.Close()

	err = gRows.Err()
	if err != nil {
		logger.Log.Error("rows error occured: ", zap.Error(err))
	}

	cRows, err := d.conn.QueryContext(ctx, `
		SELECT 
		    c.name,
		    c.counter
		FROM counter_metrics c`,
	)

	if err != nil {
		logger.Log.Error("failed to download gauge metrics")
	} else {
		for cRows.Next() {
			var m CounterMetric
			if err := cRows.Scan(&m.name, &m.value); err != nil {
				logger.Log.Error("while reading from db error occured:", zap.Error(err))
			}
			s.Counter = append(s.Counter, m)
		}
	}
	defer cRows.Close()

	err = cRows.Err()
	if err != nil {
		logger.Log.Error("rows error occured: ", zap.Error(err))
	}

	return s
}

// NewDBStorage returns new database storage
func NewDBStorage(conn *sql.DB, bootstrap bool) *DBStorage {
	db := &DBStorage{conn: conn}
	if bootstrap {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := db.Bootstrap(ctx); err != nil {
			if errors.As(err, &pgError) {
				if pgError.Code == "25P02" {
					logger.Log.Debug("rollback in bootstrap occured!")
				} else {
					logger.Log.Error("db bootstrap failed", zap.Error(err))
				}
			}
		}
	}
	return &DBStorage{conn: conn}
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
