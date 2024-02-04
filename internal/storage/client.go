package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type PGClientErrors struct {
	ErrDBConnect error
	//Transaction error
}

func NewPGClientErrors() PGClientErrors {
	return PGClientErrors{
		ErrDBConnect: errors.New("pg client couldn`t connect to server"),
	}
}

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

func (p *PGClient) GetGauge(ctx context.Context, metric *GaugeMetric) (*GaugeMetric, error) {

	sqlQuery := `SELECT id, name, gauge FROM gauge_metrics WHERE name = $1`
	tx, err := p.conn.Beginx()
	if err != nil {
		return nil, err
	}
	if err := tx.GetContext(ctx, metric, sqlQuery, metric.Name); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return metric, nil
}

func (p *PGClient) GetCounter(ctx context.Context, metric *CounterMetric) (*CounterMetric, error) {

	sqlSelect := `SELECT id, name, counter FROM counter_metrics WHERE name = $1`
	if err := p.conn.GetContext(ctx, metric, sqlSelect, metric.Name); err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return nil, err
		}
	}
	return metric, nil
}

func (p *PGClient) SetGauge(ctx context.Context, metric *GaugeMetric) (*GaugeMetric, error) {

	newValue := metric.Value
	metric, err := p.GetGauge(ctx, metric)
	if err != nil {
		return nil, err
	}
	tx, err := p.conn.Beginx()
	if err != nil {
		return nil, err
	}
	if metric.Id == 0 {
		sqlInsert := `INSERT INTO gauge_metrics (name, gauge) VALUES ($1, $2);`
		tx.MustExecContext(ctx, sqlInsert, metric.Name, newValue)
	} else {
		sqlUpdate := `UPDATE gauge_metrics SET gauge = $1 WHERE id = $2;`
		tx.MustExecContext(ctx, sqlUpdate, newValue, metric.Id)
	}
	tx.Commit()
	return metric, err
}

func (p *PGClient) SetCounter(ctx context.Context, metric *CounterMetric) (*CounterMetric, error) {

	delta := metric.Value
	metric, err := p.GetCounter(ctx, metric)
	tx, err := p.conn.Beginx()
	if err != nil {
		return nil, err
	}

	if metric.Id == 0 {
		sqlInsert := `INSERT INTO counter_metrics (name, counter) VALUES ($1, $2);`
		tx.MustExecContext(ctx, sqlInsert, metric.Name, delta)
	} else {
		sqlUpdate := `UPDATE counter_metrics SET counter = counter + $1 WHERE id = $2;`
		tx.MustExecContext(ctx, sqlUpdate, delta, metric.Id)
	}
	tx.Commit()
	return metric, err
}

func (p *PGClient) GetAllMetrics(ctx context.Context, s *StoreMetrics) (*StoreMetrics, error) {

	var allGauges []GaugeMetric
	var allCounters []CounterMetric

	sqlGaugeSelect := `SELECT id, name, gauge FROM gauge_metrics`
	tx, err := p.conn.Beginx()
	if err != nil {
		return nil, err
	}
	if err := tx.SelectContext(ctx, &allGauges, sqlGaugeSelect); err != nil {
		tx.Rollback()
		return nil, err
	}
	s.Gauge = allGauges

	sqlCounterSelect := `SELECT id, name, counter FROM counter_metrics`
	if err := tx.SelectContext(ctx, &allCounters, sqlCounterSelect); err != nil {
		tx.Rollback()
		return nil, err
	}
	s.Counter = allCounters

	tx.Commit()
	return s, nil
}
