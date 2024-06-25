package server

import (
	"context"
	"github.com/sebasttiano/Blackbird.git/internal/config"
	"github.com/sebasttiano/Blackbird.git/internal/handlers"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
)

var Serv *Server

func init() {
	repo := repository.NewMemStorage()
	srv := service.NewService(&service.Settings{Retries: 1, BackoffFactor: 1}, repo)
	views := handlers.NewServerViews(srv)
	router := views.InitRouter()
	Serv = NewServer(":3081", &views, router)
}

func TestNewServer(t *testing.T) {
	assert.IsType(t, &Server{}, Serv)
}

func TestServer_StartAndShutdown(t *testing.T) {
	cfg := &config.Config{}
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go Serv.Start(cfg)

	time.Sleep(1 * time.Second)

	// Assert socket is used
	_, err := net.Listen("tcp", ":3081")
	assert.Error(t, err)

	Serv.HandleShutdown(ctxTimeout, wg, cfg)

	wg.Wait()
	// Assert socket is free
	_, err = net.Listen("tcp", ":3081")
	assert.NoError(t, err)
}
