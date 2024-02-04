package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
)

var currentApp = newApp()

// app инкапсулирует в себя все зависимости и логику приложения
type app struct {
	store storage.Store
	views handlers.ServerViews
	conn  *sqlx.DB
}

// конструктор app
func newApp() *app {
	return &app{}
}

// Initialize принимает на вход внешние зависимости приложения и инициализирует его
func (a *app) Initialize(s *storage.StoreSettings) {

	a.conn = s.Conn

	if s.DBSave && s.Conn != nil {
		logger.Log.Info("init database storage")
		a.store = storage.NewDBStorage(a.conn, true, s.Retries, s.BackoffFactor)
	} else {
		logger.Log.Info("init mem storage")
		a.store = storage.NewMemStorage(s)
	}
	a.views = handlers.NewServerViews(a.store)
	a.views.DB = a.conn
}
