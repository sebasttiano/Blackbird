package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
)

var CurrentApp = newApp()

// app инкапсулирует в себя все зависимости и логику приложения
type app struct {
	store storage.Store
	views handlers.ServerViews
}

// конструктор app
func newApp() *app {
	return &app{}
}

// Initialize принимает на вход внешние зависимости приложения и инициализирует его
func (a *app) Initialize(s *storage.StoreSettings) {
	a.store = storage.NewMemStorage(s)
	a.views = handlers.NewServerViews(a.store)
}
