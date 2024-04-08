package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
)

var currentApp = newApp()

// app инкапсулирует в себя все зависимости и логику приложения
type app struct {
	store repository.Store
	views handlers.ServerViews
}

// конструктор app
func newApp() *app {
	return &app{}
}

// Initialize принимает на вход внешние зависимости приложения и инициализирует его
func (a *app) Initialize(s *repository.StoreSettings, key string) error {

	var err error
	if s.DBSave && s.Conn != nil {
		logger.Log.Info("init database repository")
		a.store, err = repository.NewDBStorage(s.Conn, true, s.Retries, s.BackoffFactor)
		if err != nil {
			return err
		}

	} else {
		logger.Log.Info("init mem repository")
		a.store = repository.NewMemStorage(s)
	}
	a.views = handlers.NewServerViews(a.store)
	a.views.DB = s.Conn
	a.views.SignKey = key
	return nil
}
