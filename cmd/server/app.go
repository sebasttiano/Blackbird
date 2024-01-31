package main

import (
	"database/sql"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
)

var currentApp = newApp()

// app инкапсулирует в себя все зависимости и логику приложения
type app struct {
	store storage.Store
	views handlers.ServerViews
	conn  *sql.DB
}

// конструктор app
func newApp() *app {
	return &app{}
}

// Initialize принимает на вход внешние зависимости приложения и инициализирует его
func (a *app) Initialize(s *storage.StoreSettings) {

	if s.DBSave && s.Conn != nil {
		a.conn = s.Conn
		a.store = storage.NewDBStorage(a.conn, true)
	} else {
		a.store = storage.NewMemStorage(s)
	}
	a.views = handlers.NewServerViews(a.store)
	a.views.DB = a.conn
}
