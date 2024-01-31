package storage

import "database/sql"

type StoreSettings struct {
	SyncSave     bool
	FileSave     bool
	DBSave       bool
	Conn         *sql.DB
	SaveFilePath string
}

type StoreMetrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}
