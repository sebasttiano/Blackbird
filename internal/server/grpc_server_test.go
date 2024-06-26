package server

import (
	"context"
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
)

var GServ *GRPSServer

func init() {
	repo := repository.NewMemStorage()
	srv := service.NewService(&service.Settings{Retries: 1, BackoffFactor: 1}, repo)
	GServ = NewGRPSServer(srv)
}

func TestNewGRPSServer(t *testing.T) {
	assert.IsType(t, &GRPSServer{}, GServ)
}

func TestGRPSServer_StartAndShutdown(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go GServ.Start(":4095")
	time.Sleep(1 * time.Second)

	// Assert socket is used
	_, err := net.Listen("tcp", ":4095")
	assert.Error(t, err)

	GServ.HandleShutdown(ctx, wg)
	wg.Wait()

	// Assert socket is free
	_, err = net.Listen("tcp", ":4095")
	assert.NoError(t, err)
}
