package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
)

// RetryError ошибка если все ретраи зафейлились
type RetryError struct {
	err error
}

// NewRetryError конструктор для RetryError
func NewRetryError(retries int, err error) *RetryError {
	return &RetryError{err: fmt.Errorf("function failed after %d retries. last error was %w", retries, err)}
}
func (re *RetryError) Error() string {
	return fmt.Sprintf("%v", re.err)
}

func (re *RetryError) Unwrap() error {
	return re.err
}

// pgError алиас для *pgconn.PgError
var pgError *pgconn.PgError

// ErrNoRows ошибка, если в ответе бд не вернулось ни одной строчки.
var ErrNoRows = errors.New("sql: no rows in result set")

// DBStorage тип реализующий интерфейс service.Repository
type DBStorage struct {
	conn *sqlx.DB
}

// NewDBStorage конструктор для DBStorage, с ф-цией проверки и восстановления схемы БД.
func NewDBStorage(c *sqlx.DB, bootstrap bool) (*DBStorage, error) {
	db := &DBStorage{conn: c}
	if bootstrap {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := db.Bootstrap(ctx); err != nil {
			if errors.As(err, &pgError) {
				if pgError.Code == pgerrcode.InFailedSQLTransaction {
					logger.Log.Debug("rollback in bootstrap occured!")
				} else {
					logger.Log.Error("db bootstrap failed", zap.Error(err))
				}
			}
			return nil, err
		}
	}

	return db, nil
}

// GetGauge метод из БД возвращает сохраненную метрику типа Gauge.
func (d *DBStorage) GetGauge(ctx context.Context, metric *GaugeMetric) error {
	sqlQuery := `SELECT id, name, gauge FROM gauge_metrics WHERE name = $1`

	if err := d.conn.GetContext(ctx, metric, sqlQuery, metric.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRows
		} else {
			return err
		}
	}
	return nil
}

// GetCounter метод из БД возвращает сохраненную метрику типа Counter.
func (d *DBStorage) GetCounter(ctx context.Context, metric *CounterMetric) error {
	sqlSelect := `SELECT id, name, counter FROM counter_metrics WHERE name = $1`

	if err := d.conn.GetContext(ctx, metric, sqlSelect, metric.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRows
		} else {
			return err
		}
	}
	return nil
}

// SetGauge метод сохраняет в БД метрику типа Gauge.
func (d *DBStorage) SetGauge(ctx context.Context, metric *GaugeMetric) error {
	tx, err := d.conn.Beginx()
	if err != nil {
		return err
	}

	sqlInsert := `INSERT INTO gauge_metrics (name, gauge)
                      VALUES ($1, $2)
                      ON CONFLICT (name) DO UPDATE
                      SET gauge = excluded.gauge;`

	if _, err := tx.ExecContext(ctx, sqlInsert, metric.Name, metric.Value); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// SetCounter метод сохоаняет в БД метрику типа Counter.
func (d *DBStorage) SetCounter(ctx context.Context, metric *CounterMetric) error {
	tx, err := d.conn.Beginx()
	if err != nil {
		return err
	}

	sqlInsert := `INSERT INTO counter_metrics (name, counter)
					  VALUES ($1, $2)
                      ON CONFLICT (name) DO UPDATE 
					  SET counter = counter_metrics.counter + excluded.counter;`

	if _, err := tx.ExecContext(ctx, sqlInsert, metric.Name, metric.Value); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// GetAllMetrics метод возвращает все метрики из БД
func (d *DBStorage) GetAllMetrics(ctx context.Context, sm *StoreMetrics) error {
	var allGauges []GaugeMetric
	var allCounters []CounterMetric

	sqlGaugeSelect := `SELECT id, name, gauge FROM gauge_metrics`
	if err := d.conn.SelectContext(ctx, &allGauges, sqlGaugeSelect); err != nil {
		return err
	}
	sm.Gauge = allGauges

	sqlCounterSelect := `SELECT id, name, counter FROM counter_metrics`
	if err := d.conn.SelectContext(ctx, &allCounters, sqlCounterSelect); err != nil {
		return err
	}
	sm.Counter = allCounters

	return nil
}

func (d *DBStorage) RestoreAllMetrics(gauges map[string]float64, counters map[string]int64) {}

// Bootstrap проверяет бд и создает, если надо, необходимые таблицы и типы.
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
           gauge double precision,
	       	UNIQUE(name) 
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
	       counter bigint,
	       UNIQUE(name)
	   )
	`); err != nil {
		logger.Log.Error("failed to create gauge_metrics table", zap.Error(err))
		return err
	}

	// commit
	return tx.Commit()
}
