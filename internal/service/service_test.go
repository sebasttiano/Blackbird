package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	mockservice "github.com/sebasttiano/Blackbird.git/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/sebasttiano/Blackbird.git/internal/models"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
)

func ExampleService_GetValue() {
	repo := repository.NewMemStorage()
	service := NewService(&Settings{Retries: 1, BackoffFactor: 1}, repo)
	ctx := context.TODO()

	service.SetValue(ctx, "test_gauge", "gauge", "3.33")
	value1, _ := service.GetValue(ctx, "test_gauge", "gauge")
	fmt.Println(value1)

	service.SetValue(ctx, "test_counter", "counter", "50")
	value2, _ := service.GetValue(ctx, "test_counter", "counter")
	fmt.Println(value2)

	service.SetValue(ctx, "test_counter2", "counter", "150")
	value3, _ := service.GetValue(ctx, "test_counter2", "counter")
	fmt.Println(value3)

	if err := service.SetValue(ctx, "test_counter_bad", "counter", "asd"); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 3.33
	// 50
	// 150
	// strconv.ParseInt: parsing "asd": invalid syntax
}

func ExampleService_GetModelValue() {
	repo := repository.NewMemStorage()
	service := NewService(&Settings{Retries: 1, BackoffFactor: 1}, repo)
	ctx := context.TODO()

	valueFloat := 123.456
	valueCounter := int64(100)

	m := []*models.Metrics{
		{ID: "test_gauge", MType: "gauge", Value: &valueFloat},
		{ID: "test_counter", MType: "counter", Delta: &valueCounter}}
	service.SetModelValue(ctx, m)

	m2 := models.Metrics{ID: "test_gauge", MType: "gauge"}
	m3 := models.Metrics{ID: "test_counter", MType: "counter"}

	service.GetModelValue(ctx, &m2)
	service.GetModelValue(ctx, &m3)

	fmt.Println(*m2.Value)
	fmt.Println(*m3.Delta)

	// Output:
	// 123.456
	// 100
}

func ExampleService_GetAllValues() {
	repo := repository.NewMemStorage()
	service := NewService(&Settings{Retries: 1, BackoffFactor: 1}, repo)
	ctx := context.TODO()

	valueFloat := 31.36
	valueCounter := int64(300)

	m := []*models.Metrics{
		{ID: "test_gauge", MType: "gauge", Value: &valueFloat},
		{ID: "test_counter", MType: "counter", Delta: &valueCounter}}

	service.SetModelValue(ctx, m)

	sm := service.GetAllValues(ctx)
	for _, v := range sm.Counter {
		fmt.Println(v.Value)
	}
	for _, v := range sm.Gauge {
		fmt.Println(v.Value)
	}

	// Output:
	// 300
	// 31.36
}

func TestService_Save(t *testing.T) {

	testTable := []struct {
		name string
		repo string
		sm   repository.StoreMetrics
		path string
		err  error
	}{
		{
			name: "OK save MemStorage",
			repo: "memory",
			sm: repository.StoreMetrics{
				Gauge:   make([]repository.GaugeMetric, 0),
				Counter: make([]repository.CounterMetric, 0),
			},
			path: "/tmp/metrics-db.json",
			err:  nil,
		},
		{
			name: "NOT OK save. DBStorage",
			repo: "db",
			sm: repository.StoreMetrics{
				Gauge:   make([]repository.GaugeMetric, 0),
				Counter: make([]repository.CounterMetric, 0),
			},
			path: "/tmp/metrics-db.json",
			err:  ErrNotSupported,
		},
		{
			name: "NOT OK save. Another",
			repo: "mock",
			sm: repository.StoreMetrics{
				Gauge:   make([]repository.GaugeMetric, 0),
				Counter: make([]repository.CounterMetric, 0),
			},
			path: "/tmp/metrics-db.json",
			err:  ErrNotSupported,
		},
		{
			name: "NOT OK save. Invalid Path",
			repo: "memory",
			sm: repository.StoreMetrics{
				Gauge:   make([]repository.GaugeMetric, 0),
				Counter: make([]repository.CounterMetric, 0),
			},
			path: "",
			err:  errors.New("can`t save to file. no file path specify"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			var repo Repository
			switch tt.repo {
			case "memory":
				repo = repository.NewMemStorage()
			case "db":
				sql := &sqlx.DB{}
				repo, _ = repository.NewDBStorage(sql, false)
			default:
				c := gomock.NewController(t)
				defer c.Finish()
				repo = mockservice.NewMockRepository(c)
			}

			service := NewService(&Settings{Retries: 1, BackoffFactor: 1, SaveFilePath: tt.path}, repo)
			err := service.Save()
			if tt.err != nil {
				if assert.Errorf(t, err, err.Error()) {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Restore(t *testing.T) {

	testTable := []struct {
		name string
		repo string
		path string
		err  error
	}{
		{
			name: "OK restore MemStorage",
			repo: "memory",
			path: "/tmp/metrics-db.json",
			err:  nil,
		},
		{
			name: "NOT OK restore DBstorage",
			repo: "db",
			path: "",
			err:  ErrNotSupported,
		},
		{
			name: "NOT OK restore default",
			repo: "mock",
			path: "",
			err:  ErrNotSupported,
		},
		{
			name: "NOT OK restore. Invalid Path",
			repo: "memory",
			path: "",
			err:  errors.New("can`t restore from file. no file path specify"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			var repo Repository
			switch tt.repo {
			case "memory":
				repo = repository.NewMemStorage()
			case "db":
				sql := &sqlx.DB{}
				repo, _ = repository.NewDBStorage(sql, false)
			default:
				c := gomock.NewController(t)
				defer c.Finish()
				repo = mockservice.NewMockRepository(c)
			}

			service := NewService(&Settings{Retries: 1, BackoffFactor: 1, SaveFilePath: tt.path}, repo)
			err := service.Restore()
			if tt.err != nil {
				if assert.Errorf(t, err, err.Error()) {
					assert.Equal(t, tt.err, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestService_RetryDBError(t *testing.T) {

	testError := NewRetryDBError(3, errors.New("failed to connect to database"))
	assert.Equal(t, testError.Error(), "function failed after 3 retries. last error was failed to connect to database")
	assert.Equal(t, testError.Unwrap().Error(), "failed to connect to database")
}
