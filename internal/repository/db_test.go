package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestDBStorage_GetGauge(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s, err := NewDBStorage(db, false)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when create db storage type", err)
	}
	testTable := []struct {
		name string
		s    *DBStorage
		m    *GaugeMetric
		mock func()
		want *GaugeMetric
		err  error
	}{
		{
			name: "OK",
			s:    s,
			m:    &GaugeMetric{ID: 0, Name: "test_gauge", Value: 0},
			mock: func() {
				selectMockGaugeRows := sqlxmock.NewRows([]string{"id", "name", "gauge"}).AddRow(1, "test_gauge", 137.3)
				mock.ExpectQuery("SELECT id, name, gauge FROM").WithArgs("test_gauge").WillReturnRows(selectMockGaugeRows)
			},
			want: &GaugeMetric{ID: 1, Name: "test_gauge", Value: 137.3},
			err:  nil,
		},
		{
			name: "NOT OK. ErrNoRows",
			s:    s,
			m:    &GaugeMetric{ID: 0, Name: "test_gauge", Value: 0},
			mock: func() {
				mock.ExpectQuery("SELECT id, name, gauge FROM").WithArgs("test_gauge").WillReturnError(sql.ErrNoRows)
			},
			err: ErrNoRows,
		},
		{
			name: "NOT OK. something went wrong",
			s:    s,
			m:    &GaugeMetric{ID: 0, Name: "test_gauge", Value: 0},
			mock: func() {
				mock.ExpectQuery("SELECT id, name, gauge FROM").WithArgs("test_gauge").WillReturnError(errors.New("something went wrong"))
			},
			err: errors.New("something went wrong"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.s.GetGauge(context.TODO(), tt.m)
			if tt.err != nil {
				if assert.Errorf(t, err, err.Error()) {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.m)
			}
		})
	}

}

func TestDBStorage_GetCounter(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s, err := NewDBStorage(db, false)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when create db storage type", err)
	}
	testTable := []struct {
		name string
		s    *DBStorage
		m    *CounterMetric
		mock func()
		want *CounterMetric
		err  error
	}{
		{
			name: "OK",
			s:    s,
			m:    &CounterMetric{ID: 0, Name: "test_counter", Value: 0},
			mock: func() {
				selectMockCounterRows := sqlxmock.NewRows([]string{"id", "name", "counter"}).AddRow(1, "test_counter", 250)
				mock.ExpectQuery("SELECT id, name, counter FROM").WithArgs("test_counter").WillReturnRows(selectMockCounterRows)
			},
			want: &CounterMetric{ID: 1, Name: "test_counter", Value: 250},
			err:  nil,
		},
		{
			name: "NOT OK. ErrNoRows",
			s:    s,
			m:    &CounterMetric{ID: 0, Name: "test_counter", Value: 0},
			mock: func() {
				mock.ExpectQuery("SELECT id, name, counter FROM").WithArgs("test_counter").WillReturnError(sql.ErrNoRows)
			},
			err: ErrNoRows,
		},
		{
			name: "NOT OK. something went wrong",
			s:    s,
			m:    &CounterMetric{ID: 0, Name: "test_counter", Value: 0},
			mock: func() {
				mock.ExpectQuery("SELECT id, name, counter FROM").WithArgs("test_counter").WillReturnError(errors.New("something went wrong"))
			},
			err: errors.New("something went wrong"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.s.GetCounter(context.TODO(), tt.m)
			if tt.err != nil {
				if assert.Errorf(t, err, err.Error()) {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.m)
			}
		})
	}
}

func TestDBStorage_SetGauge(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s, err := NewDBStorage(db, false)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when create db storage type", err)
	}

	testTable := []struct {
		name string
		s    *DBStorage
		m    *GaugeMetric
		mock func()
		err  error
	}{
		{
			name: "OK",
			s:    s,
			m:    &GaugeMetric{ID: 0, Name: "test_gauge", Value: 338.1},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs("test_gauge", 338.1).WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			err: nil,
		},
		{
			name: "NOT OK",
			s:    s,
			m:    &GaugeMetric{ID: 0, Name: "test_gauge", Value: 338.1},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs("test_gauge", 338.1).WillReturnError(errors.New("something went wrong"))
				mock.ExpectRollback()
			},
			err: errors.New("something went wrong"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.s.SetGauge(context.TODO(), tt.m)
			if tt.err != nil {
				if assert.Errorf(t, err, tt.err.Error()) {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDBStorage_SetCounter(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s, err := NewDBStorage(db, false)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when create db storage type", err)
	}

	testTable := []struct {
		name string
		s    *DBStorage
		m    *CounterMetric
		mock func()
		err  error
	}{
		{
			name: "OK",
			s:    s,
			m:    &CounterMetric{ID: 0, Name: "test_counter", Value: 113},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs("test_counter", 113).WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			err: nil,
		},
		{
			name: "NOT OK",
			s:    s,
			m:    &CounterMetric{ID: 0, Name: "test_counter", Value: 113},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WithArgs("test_counter", 113).WillReturnError(errors.New("something went wrong"))
				mock.ExpectRollback()
			},
			err: errors.New("something went wrong"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.s.SetCounter(context.TODO(), tt.m)
			if tt.err != nil {
				if assert.Errorf(t, err, tt.err.Error()) {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDBStorage_GetAllMetrics(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s, err := NewDBStorage(db, false)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when create db storage type", err)
	}

	testTable := []struct {
		name string
		s    *DBStorage
		sm   *StoreMetrics
		mock func()
		want *StoreMetrics
		err  error
	}{
		{
			name: "OK",
			s:    s,
			sm:   &StoreMetrics{make([]GaugeMetric, 0, 2), make([]CounterMetric, 0, 2)},
			mock: func() {
				selectMockGaugeRows := sqlxmock.NewRows([]string{"id", "name", "gauge"}).AddRow(1, "test_gauge1", 338.1).AddRow(2, "test_gauge2", 187.3)
				mock.ExpectQuery("SELECT id, name, gauge FROM gauge_metrics").WillReturnRows(selectMockGaugeRows)
				selectMockCounterRows := sqlxmock.NewRows([]string{"id", "name", "counter"}).AddRow(1, "test_counter1", 777).AddRow(2, "test_counter2", 90)
				mock.ExpectQuery("SELECT id, name, counter FROM counter_metrics").WillReturnRows(selectMockCounterRows)
			},
			want: &StoreMetrics{
				Gauge:   []GaugeMetric{{1, "test_gauge1", 338.1}, {2, "test_gauge2", 187.3}},
				Counter: []CounterMetric{{1, "test_counter1", 777}, {2, "test_counter2", 90}},
			},
			err: nil,
		},
		{
			name: "NOT OK. gauge select failed",
			s:    s,
			sm:   &StoreMetrics{make([]GaugeMetric, 0, 2), make([]CounterMetric, 0, 2)},
			mock: func() {
				mock.ExpectQuery("SELECT id, name, gauge FROM gauge_metrics").WillReturnError(errors.New("something went wrong with gauge select"))
			},
			err: errors.New("something went wrong with gauge select"),
		},
		{
			name: "NOT OK. counter select failed",
			s:    s,
			sm:   &StoreMetrics{make([]GaugeMetric, 0, 2), make([]CounterMetric, 0, 2)},
			mock: func() {
				selectMockGaugeRows := sqlxmock.NewRows([]string{"id", "name", "gauge"}).AddRow(1, "test_gauge1", 338.1).AddRow(2, "test_gauge2", 187.3)
				mock.ExpectQuery("SELECT id, name, gauge FROM gauge_metrics").WillReturnRows(selectMockGaugeRows)
				mock.ExpectQuery("SELECT id, name, counter FROM counter_metrics").WillReturnError(errors.New("something went wrong with counter select"))
			},
			err: errors.New("something went wrong with counter select"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.s.GetAllMetrics(context.TODO(), tt.sm)
			if tt.err != nil {
				if assert.Errorf(t, err, tt.err.Error()) {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.sm)
			}
		})
	}
}

func TestDBStorage_Bootstrap(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s, err := NewDBStorage(db, false)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when create db storage type", err)
	}

	testTable := []struct {
		name string
		s    *DBStorage
		mock func()
		err  error
	}{
		{
			name: "OK",
			s:    s,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS gauge_metrics").WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS counter_metrics").WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			err: nil,
		},
		{
			name: "NOT OK. create gauge table failed",
			s:    s,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS gauge_metrics").WillReturnError(errors.New("failed to create gauge_metrics table"))
			},
			err: errors.New("failed to create gauge_metrics table"),
		},
		{
			name: "NOT OK. create counter table failed",
			s:    s,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS gauge_metrics").WillReturnResult(sqlxmock.NewResult(1, 1))
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS counter_metrics").WillReturnError(errors.New("failed to create counter_metrics table"))
			},
			err: errors.New("failed to create counter_metrics table"),
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.s.Bootstrap(context.TODO())
			if tt.err != nil {
				if assert.Errorf(t, err, tt.err.Error()) {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
