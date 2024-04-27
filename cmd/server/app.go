package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/sebasttiano/Blackbird.git/internal/service"
)

var currentApp = newApp()

// app инкапсулирует в себя все зависимости и логику приложения
type app struct {
	service *service.Service
	views   handlers.ServerViews
}

// newApp конструктор app
func newApp() *app {
	return &app{}
}

// Initialize принимает на вход внешние зависимости приложения и инициализирует его
func (a *app) Initialize(s *service.ServiceSettings, key string) error {
	var err error
	var repo service.Repository

	if s.DBSave && s.Conn != nil {
		logger.Log.Info("init database repository")
		repo, err = repository.NewDBStorage(s.Conn, true)
		if err != nil {
			return err
		}
	} else {
		logger.Log.Info("init mem repository")
		repo = repository.NewMemStorage()
	}

	a.service = service.NewService(s, repo)
	a.views = handlers.NewServerViews(a.service)
	a.views.DB = s.Conn
	a.views.SignKey = key
	return nil
}
