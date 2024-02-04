package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	"time"
)

type PGClientErrors struct {
	ErrDBConnect error
	ErrNoRows    error
}

func NewPGClientErrors() PGClientErrors {
	return PGClientErrors{
		ErrDBConnect: errors.New("pg client couldn`t connect to server"),
		ErrNoRows:    sql.ErrNoRows,
	}
}

// PGClient define postgres client
type PGClient struct {
	conn        *sqlx.DB
	retryDelays []uint
	Errors      PGClientErrors
}

// NewPGClient creates client and retries chrony calculates
func NewPGClient(c *sqlx.DB, retries uint, backoffFactor uint) *PGClient {
	var ri []uint
	for i := 1; i <= int(retries); i++ {
		ri = append(ri, backoffFactor*uint(i)-1)
	}
	return &PGClient{
		conn:        c,
		retryDelays: ri,
		Errors:      NewPGClientErrors(),
	}
}

// GetGauge method to get from gauge_metrics table
func (p *PGClient) GetGauge(ctx context.Context, metric *GaugeMetric) (*GaugeMetric, error) {

	sqlQuery := `SELECT id, name, gauge FROM gauge_metrics WHERE name = $1`

	if err := p.conn.GetContext(ctx, metric, sqlQuery, metric.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return metric, p.Errors.ErrNoRows
		} else {
			return nil, err
		}
	}
	return metric, nil
}

// GetCounter method to get from counter_metrics table
func (p *PGClient) GetCounter(ctx context.Context, metric *CounterMetric) (*CounterMetric, error) {

	sqlSelect := `SELECT id, name, counter FROM counter_metrics WHERE name = $1`

	if err := p.conn.GetContext(ctx, metric, sqlSelect, metric.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return metric, p.Errors.ErrNoRows
		} else {
			return nil, err
		}
	}
	return metric, nil
}

// SetGauge method inserts or updates rows in gauge_metrics table
func (p *PGClient) SetGauge(ctx context.Context, metric *GaugeMetric) (*GaugeMetric, error) {

	newValue := metric.Value
	metric, err := p.GetGauge(ctx, metric)
	if err != nil {
		if !errors.Is(err, p.Errors.ErrNoRows) {
			return nil, err
		}
	}
	tx, err := p.conn.Beginx()
	if err != nil {
		return nil, err
	}
	if metric.ID == 0 {
		sqlInsert := `INSERT INTO gauge_metrics (name, gauge) VALUES ($1, $2);`
		tx.MustExecContext(ctx, sqlInsert, metric.Name, newValue)
	} else {
		sqlUpdate := `UPDATE gauge_metrics SET gauge = $1 WHERE id = $2;`
		tx.MustExecContext(ctx, sqlUpdate, newValue, metric.ID)
	}
	tx.Commit()
	return metric, err
}

// SetCounter method inserts or updates rows in counter_metrics table
func (p *PGClient) SetCounter(ctx context.Context, metric *CounterMetric) (*CounterMetric, error) {

	delta := metric.Value
	metric, err := p.GetCounter(ctx, metric)
	if err != nil {
		if !errors.Is(err, p.Errors.ErrNoRows) {
			return nil, err
		}
	}
	tx, err := p.conn.Beginx()
	if err != nil {
		return nil, err
	}

	if metric.ID == 0 {
		sqlInsert := `INSERT INTO counter_metrics (name, counter) VALUES ($1, $2);`
		tx.MustExecContext(ctx, sqlInsert, metric.Name, delta)
	} else {
		sqlUpdate := `UPDATE counter_metrics SET counter = counter + $1 WHERE id = $2;`
		tx.MustExecContext(ctx, sqlUpdate, delta, metric.ID)
	}
	tx.Commit()
	return metric, err
}

// GetAllMetrics takes all rows from gauge_metrics and counter_metrics tables
func (p *PGClient) GetAllMetrics(ctx context.Context, s *StoreMetrics) (*StoreMetrics, error) {

	var allGauges []GaugeMetric
	var allCounters []CounterMetric

	sqlGaugeSelect := `SELECT id, name, gauge FROM gauge_metrics`
	if err := p.conn.SelectContext(ctx, &allGauges, sqlGaugeSelect); err != nil {
		return nil, err
	}
	s.Gauge = allGauges

	sqlCounterSelect := `SELECT id, name, counter FROM counter_metrics`
	if err := p.conn.SelectContext(ctx, &allCounters, sqlCounterSelect); err != nil {
		return nil, err
	}
	s.Counter = allCounters

	return s, nil
}

// RetryError custom error for retry
type RetryError struct {
	Err error
}

func (re *RetryError) Error() string {
	return fmt.Sprintf("%v", re.Err)
}

func (re *RetryError) Unwrap() error {
	return re.Err
}

// Retry function repeat failed queries to db based within retry delays, ignores sql.ErrNoRows
func Retry[T any](ctx context.Context, retryDelays []uint, f func(context.Context, T) (T, error), arg T) (T, error) {

	var retries = len(retryDelays)
	for _, delay := range retryDelays {
		select {
		case <-ctx.Done():
			return arg, ctx.Err()
		default:
			result, err := f(ctx, arg)
			retries -= 1
			if err != nil {
				pgErr := NewPGClientErrors()
				if !errors.Is(err, pgErr.ErrNoRows) {
					logger.Log.Error(fmt.Sprintf("Request to server failed. retrying in %d seconds... Retries left %d\n", delay, retries), zap.Error(err))
					time.Sleep(time.Duration(delay) * time.Second)
					if retries == 0 {
						return arg, &RetryError{Err: err}
					}
				} else {
					return arg, err
				}
			} else {
				return result, nil
			}
		}
	}
	return arg, nil
}
