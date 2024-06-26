package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"go.uber.org/zap"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

// ErrNotSupported общая ошибка, если сервис не поддерживает действие.
var ErrNotSupported = errors.New("service not supported")
var ErrUnknownMetricType = errors.New("unknown metric type. only gauge and counter are available")

// RetryDBError тип реализующий интерфейс Error, записывает количество ретраев и заворачивает ошибку ф-ция.
type RetryDBError struct {
	Retries int
	Err     error
}

// Error метод интерфейса записывает количество ретраев и оригинальную ошибку.
func (e RetryDBError) Error() string {
	return fmt.Sprintf("function failed after %d retries. last error was %v", e.Retries, e.Err)
}

// Unwrap метод интерфейса возвращает упакованную оригинальную ошибку.
func (e RetryDBError) Unwrap() error {
	return e.Err
}

// NewRetryDBError конструктор для  RetryDBError
func NewRetryDBError(retries int, err error) *RetryDBError {
	return &RetryDBError{retries, err}
}

// Settings настройки сервиса.
type Settings struct {
	SyncSave      bool
	FileSave      bool
	DBSave        bool
	Conn          *sqlx.DB
	SaveFilePath  string
	Retries       uint
	BackoffFactor uint
	TrustedSubnet *net.IPNet
}

// Service реализует интерфейс MetricService.
type Service struct {
	Settings     *Settings
	fileRestorer FileService
	repo         Repository
	retries      []uint
}

// NewService конструктор для Service.
func NewService(serviceSettings *Settings, repo Repository) *Service {
	var ri []uint
	for i := 1; i <= int(serviceSettings.Retries); i++ {
		ri = append(ri, serviceSettings.BackoffFactor*uint(i)-1)
	}
	return &Service{serviceSettings, NewFileHanlder(serviceSettings.SaveFilePath), repo, ri}
}

// MetricService интерфейс описывающий работу с метриками
type MetricService interface {
	GetValue(ctx context.Context, string, metricType string) (interface{}, error)
	GetModelValue(ctx context.Context, metric *models.Metrics) error
	SetValue(ctx context.Context, metricName string, metricType string, metricValue string) error
	SetModelValue(ctx context.Context, metrics []*models.Metrics) error
	GetAllValues(ctx context.Context) *repository.StoreMetrics
	Save() error
	Restore() error
}

// Repository интерфейс описывающий сохранение и чтение метрик из хранилища.
type Repository interface {
	GetGauge(ctx context.Context, metric *repository.GaugeMetric) error
	GetCounter(ctx context.Context, metric *repository.CounterMetric) error
	SetGauge(ctx context.Context, metric *repository.GaugeMetric) error
	SetCounter(ctx context.Context, metric *repository.CounterMetric) error
	GetAllMetrics(ctx context.Context, s *repository.StoreMetrics) error
	RestoreAllMetrics(gauges map[string]float64, counters map[string]int64)
}

// GetValue возвращает или Gauge, или Counter метрики.
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
		return m.Value, nil
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
		return m.Value, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownMetricType, metricType)
	}
}

// GetModelValue маппит данные из хранилища в структуру.
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

// SetValue сохраняет или Gauge, или Counter метрики.
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
	default:
		return fmt.Errorf("%w: %s", ErrUnknownMetricType, metricType)
	}

	if s.Settings.SyncSave {
		if err := s.Save(); err != nil {
			logger.Log.Error("couldn`t save to the file", zap.Error(err))
			return err
		}
	}
	return nil
}

// SetModelValue сохраняет или Gauge, или Counter метрики из моделек.
func (s *Service) SetModelValue(ctx context.Context, metrics []*models.Metrics) error {
	for _, metric := range metrics {
		if metric.ID == "" {
			return errors.New("name of the metric is required")
		}

		switch metric.MType {
		case "gauge":
			if metric.Value == nil {
				return fmt.Errorf("value of the gauge is required. %s", metric.ID)
			}
			if err := s.SetValue(ctx, metric.ID, metric.MType, fmt.Sprintf("%.12f", *metric.Value)); err != nil {
				return err
			}

		case "counter":
			if metric.Delta == nil {
				return fmt.Errorf("value of the counter is required. %s", metric.ID)
			}
			if err := s.SetValue(ctx, metric.ID, metric.MType, fmt.Sprintf("%d", *metric.Delta)); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: %s", ErrUnknownMetricType, metric.MType)
		}
	}
	return nil
}

// GetAllValues забирает все метрики из хранилища.
func (s *Service) GetAllValues(ctx context.Context) (sm *repository.StoreMetrics) {
	sm = &repository.StoreMetrics{Gauge: make([]repository.GaugeMetric, 0), Counter: make([]repository.CounterMetric, 0)}

	s.Retry(ctx, s.retries, func(ctx context.Context) error {
		return s.repo.GetAllMetrics(ctx, sm)
	})
	return sm
}

// Save сохраняет в хранилище, если оно типа repository.MemStorage.
func (s *Service) Save() error {
	switch s.repo.(type) {
	case *repository.MemStorage:
		var sm repository.StoreMetrics

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.repo.GetAllMetrics(ctx, &sm); err != nil {
			return fmt.Errorf("failed to save metrics, %w", err)
		}

		gauges := make(map[string]float64)
		counters := make(map[string]int64)

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

// Restore восстанавливает их хранилища, если оно типа repository.MemStorage.
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

// Retry метод повтора функций с задержками при повторных попытках, игнорирует sql.ErrNoRows
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
						return NewRetryDBError(retries, err)
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
