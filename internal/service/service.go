package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var ErrNotSupported = errors.New("service not supported")

type ErrRetryDB struct {
	Retries int
	Err     error
}

func (e ErrRetryDB) Error() string {
	return fmt.Sprintf("function failed after %d retries. last error was %v", e.Retries, e.Err)
}

func (e ErrRetryDB) Unwrap() error {
	return e.Err
}

func NewErrRetryDB(retries int, err error) *ErrRetryDB {
	return &ErrRetryDB{retries, err}
}

type ServiceSettings struct {
	SyncSave      bool
	FileSave      bool
	DBSave        bool
	Conn          *sqlx.DB
	SaveFilePath  string
	Retries       uint
	BackoffFactor uint
}

type Service struct {
	Settings     *ServiceSettings
	fileRestorer FileService
	repo         Repository
	retries      []uint
}

func NewService(serviceSettings *ServiceSettings, repo Repository) *Service {

	var ri []uint
	for i := 1; i <= int(serviceSettings.Retries); i++ {
		ri = append(ri, serviceSettings.BackoffFactor*uint(i)-1)
	}
	fmt.Println(ri)
	return &Service{serviceSettings, NewFileHanlder(serviceSettings.SaveFilePath), repo, ri}
}

type MetricService interface {
	GetValue(ctx context.Context, string, metricType string) (interface{}, error)
	GetModelValue(ctx context.Context, metric *models.Metrics) error
	SetValue(ctx context.Context, metricName string, metricType string, metricValue string) error
	SetModelValue(ctx context.Context, metrics []*models.Metrics) error
	GetAllValues(ctx context.Context) *repository.StoreMetrics
	Save() error
	Restore() error
}

type Repository interface {
	GetGauge(ctx context.Context, metric *repository.GaugeMetric) error
	GetCounter(ctx context.Context, metric *repository.CounterMetric) error
	SetGauge(ctx context.Context, metric *repository.GaugeMetric) error
	SetCounter(ctx context.Context, metric *repository.CounterMetric) error
	GetAllMetrics(ctx context.Context, s *repository.StoreMetrics) error
	RestoreAllMetrics(gauges map[string]float64, counters map[string]int64)
}

// GetValue returns either gauge or counter metrics
func (s *Service) GetValue(ctx context.Context, metricName string, metricType string) (interface{}, error) {

	switch metricType {
	case "gauge":
		m := repository.GaugeMetric{Name: metricName}
		var err error
		err = s.Retry(ctx, s.retries, func(ctx context.Context) error {
			return s.repo.GetGauge(ctx, &m)
		},
		)
		if err != nil && !errors.Is(repository.ErrNoRows, err) {
			return nil, fmt.Errorf("failed to load gauge metric %w", err)
		}
		return m, nil
	case "counter":
		m := repository.CounterMetric{Name: metricName}
		var err error
		err = s.Retry(ctx, s.retries, func(ctx context.Context) error {
			return s.repo.GetCounter(ctx, &m)
		},
		)
		if err != nil && !errors.Is(repository.ErrNoRows, err) {
			return nil, fmt.Errorf("failed to load gauge metric %w", err)
		}
		return m, nil
	default:
		return nil, errors.New("error: unknown metric type. only gauge and counter are available")
	}
}

func (s *Service) GetModelValue(ctx context.Context, metric *models.Metrics) error {

	if metric.ID == "" {
		return errors.New("name of the metric is required")
	}

	value, err := s.GetValue(ctx, metric.ID, metric.MType)
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
func (s *Service) SetValue(ctx context.Context, metricName string, metricType string, metricValue string) error {

	switch metricType {
	case "gauge":
		valueFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return err
		}
		m := repository.GaugeMetric{Name: metricName, Value: valueFloat}
		err = s.Retry(ctx, s.retries, func(ctx context.Context) error {
			return s.repo.SetGauge(ctx, &m)
		})
		if err != nil {
			return err
		}
	case "counter":
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return err
		}
		m := repository.CounterMetric{Name: metricName, Value: intValue}
		err = s.Retry(ctx, s.retries, func(ctx context.Context) error {
			return s.repo.SetCounter(ctx, &m)
		})
		if err != nil {
			return err
		}
	}

	if s.Settings.SyncSave {
		if err := s.Save(); err != nil {
			logger.Log.Error("couldn`t save to the file", zap.Error(err))
			return err
		}
	}
	return nil
}

// SetModelValue saves either gauge or counter metrics from model
func (s *Service) SetModelValue(ctx context.Context, metrics []*models.Metrics) error {

	for _, metric := range metrics {

		if metric.ID == "" {
			return errors.New("name of the metric is required")
		}

		switch metric.MType {
		case "gauge":
			if metric.Value == nil {
				return errors.New("value of the gauge is required")
			}
			if err := s.SetValue(ctx, metric.ID, metric.MType, fmt.Sprintf("%.12f", *metric.Value)); err != nil {
				return err
			}

		case "counter":
			if metric.Delta == nil {
				return errors.New("value of the counter is required")
			}
			if err := s.SetValue(ctx, metric.ID, metric.MType, fmt.Sprintf("%d", *metric.Delta)); err != nil {
				return err
			}
		default:
			return errors.New("error: unknown metric type. Only gauge and counter are available")
		}
	}
	return nil
}

// GetAllValues get all metrics from db and returns in raw format
func (s *Service) GetAllValues(ctx context.Context) (sm *repository.StoreMetrics) {

	sm = &repository.StoreMetrics{Gauge: make([]repository.GaugeMetric, 0), Counter: make([]repository.CounterMetric, 0)}

	s.Retry(ctx, s.retries, func(ctx context.Context) error {
		return s.repo.GetAllMetrics(ctx, sm)
	})
	return sm
}

func (s *Service) Save() error {
	switch s.repo.(type) {
	case *repository.MemStorage:
		var sm repository.StoreMetrics

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.repo.GetAllMetrics(ctx, &sm); err != nil {
			return fmt.Errorf("failed to save metrics, %w", err)
		}

		var gauges map[string]float64
		var counters map[string]int64

		for _, metric := range sm.Gauge {
			gauges[metric.Name] = metric.Value
		}

		for _, metric := range sm.Counter {
			counters[metric.Name] = metric.Value
		}

		return s.fileRestorer.Save(gauges, counters)
	default:
		return ErrNotSupported
	}
}

func (s *Service) Restore() error {
	switch s.repo.(type) {
	case *repository.MemStorage:
		gauges, counters, err := s.fileRestorer.Restore()
		if err != nil {
			return err
		}
		s.repo.RestoreAllMetrics(gauges, counters)
		return nil
	default:
		return ErrNotSupported
	}
}

// Retry method repeat functions calls within retry delays, ignores sql.ErrNoRows
func (s *Service) Retry(ctx context.Context, retryDelays []uint, f func(ctx context.Context) error) error {

	var retries = len(retryDelays)
	for _, delay := range retryDelays {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := f(ctx)
			retries -= 1
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					logger.Log.Error(fmt.Sprintf("Request to server failed. retrying in %d seconds... Retries left %d\n", delay, retries), zap.Error(err))
					time.Sleep(time.Duration(delay) * time.Second)
					if retries == 0 {
						return NewErrRetryDB(retries, err)
					}
				} else {
					return err
				}
			} else {
				return nil
			}
		}
	}
	return nil
}
