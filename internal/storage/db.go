package storage

import (
	"context"
	"database/sql"
	"errors"
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

//	type DBStorageErrors struct {
//		ErrConnect error
//	}
var pgError *pgconn.PgError

func (d *DBStorage) GetValue(metricName string, metricType string) (interface{}, error) {
	return nil, nil
}

func (d *DBStorage) GetModelValue(metrics *models.Metrics) error {
	return nil
}

func (d *DBStorage) SetValue(metricName string, metricType string, metricValue string) error {
	return nil
}

func (d *DBStorage) SetModelValue(metric *models.Metrics) error {
	return nil
}

func (d *DBStorage) Save() error {
	return nil
}

func (d *DBStorage) Restore() error {
	return nil
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
            id varchar(128) PRIMARY KEY,
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
	       id varchar(128) PRIMARY KEY,
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
