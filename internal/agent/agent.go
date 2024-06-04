// Package agent содержит основные структуры, типы и методы для работы компонента программы - агент.
package agent

import (
	"context"
	"errors"
	"math/rand"
	"regexp"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
)

var ErrInitSender = errors.New("failed to init sender")
var ErrSendToRepo = errors.New("failed to send to repo")

// MetricsSet - структура в которой перечислены основные runtime метрики приложения.
type MetricsSet struct {
	Alloc,
	TotalAlloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	NumForcedGC,
	NumGC,
	OtherSys,
	PauseTotalNs,
	StackInuse,
	StackSys,
	Sys,
	RandomValue float64
	PollCount int64
}

// GopsutilMetricsSet - структура, которая содержит общее потребление ресурсов на хосте.
type GopsutilMetricsSet struct {
	TotalMemory,
	FreeMemory,
	CPUUtilization float64
}

type Sender interface {
	SendToRepo(jobsMetrics <-chan MetricsSet, jobsGMetrics <-chan GopsutilMetricsSet) error
}

// Agent - тип, который реализует сущность агент.
type Agent struct {
	getCounter int64
	rtm        runtime.MemStats
	Metrics    MetricsSet
	GMetrics   GopsutilMetricsSet
	WG         sync.WaitGroup
	Sender     Sender
}

// NewAgent - конструктор для типа Agent.
func NewAgent(serverAddr string, clientRetries int, backoffFactor uint, signKey string, publicKey []byte, grpcServer string) (*Agent, error) {
	getCounter := new(int64)
	re, _ := regexp.Compile("^.+://(.+$)")
	addr := re.FindAllStringSubmatch(serverAddr, 1)
	x, err := common.GetLocalIP(addr[0][1])
	var xRealIP string
	if err != nil {
		logger.Log.Warn("failed to get local IP", zap.String("serverAddr", serverAddr))
		xRealIP = ""
	} else {
		xRealIP = x.String()
	}

	if grpcServer != "" {
		gClient, err := NewGRPCClient(grpcServer)
		if err != nil {
			return nil, err
		}
		return &Agent{
			getCounter: *getCounter,
			Sender:     gClient,
		}, nil
	}
	return &Agent{
		getCounter: *getCounter,
		Sender: &HTTPSender{
			client:    common.NewHTTPClient(serverAddr, clientRetries, backoffFactor),
			signKey:   signKey,
			publicKey: common.UnmarshalRSAPublic(publicKey),
			XRealIP:   xRealIP,
		},
	}, nil
}

// GetMetrics - метод для сбора runtime метрик приложения с определенным интервалом
func (a *Agent) GetMetrics(ctx context.Context, getInterval time.Duration, jobs chan<- MetricsSet) {
	defer close(jobs)
	tick := time.NewTicker(getInterval)

	for range tick.C {
		atomic.AddInt64(&a.getCounter, 1)
		runtime.ReadMemStats(&a.rtm)

		logger.Log.Info("collect memstats successfully")
		a.Metrics.Alloc = float64(a.rtm.Alloc)
		a.Metrics.TotalAlloc = float64(a.rtm.TotalAlloc)
		a.Metrics.BuckHashSys = float64(a.rtm.BuckHashSys)
		a.Metrics.Frees = float64(a.rtm.Frees)
		a.Metrics.GCCPUFraction = a.rtm.GCCPUFraction
		a.Metrics.GCSys = float64(a.rtm.GCSys)
		a.Metrics.HeapAlloc = float64(a.rtm.HeapAlloc)
		a.Metrics.HeapIdle = float64(a.rtm.HeapIdle)
		a.Metrics.HeapInuse = float64(a.rtm.HeapInuse)
		a.Metrics.HeapObjects = float64(a.rtm.HeapObjects)
		a.Metrics.HeapReleased = float64(a.rtm.HeapReleased)
		a.Metrics.HeapSys = float64(a.rtm.HeapSys)
		a.Metrics.LastGC = float64(a.rtm.LastGC)
		a.Metrics.Lookups = float64(a.rtm.Lookups)
		a.Metrics.MCacheInuse = float64(a.rtm.MCacheInuse)
		a.Metrics.MCacheSys = float64(a.rtm.MCacheSys)
		a.Metrics.MSpanInuse = float64(a.rtm.MSpanInuse)
		a.Metrics.MSpanSys = float64(a.rtm.MSpanSys)
		a.Metrics.Mallocs = float64(a.rtm.Mallocs)
		a.Metrics.NextGC = float64(a.rtm.NextGC)
		a.Metrics.NumForcedGC = float64(a.rtm.NumForcedGC)
		a.Metrics.NumGC = float64(a.rtm.NumGC)
		a.Metrics.OtherSys = float64(a.rtm.OtherSys)
		a.Metrics.PauseTotalNs = float64(a.rtm.PauseTotalNs)
		a.Metrics.StackInuse = float64(a.rtm.StackInuse)
		a.Metrics.StackSys = float64(a.rtm.StackSys)
		a.Metrics.Sys = float64(a.rtm.Sys)
		a.Metrics.PollCount = a.getCounter
		a.Metrics.RandomValue = rand.Float64()

		select {
		case jobs <- a.Metrics:
		case <-ctx.Done():
			tick.Stop()
			a.WG.Done()
			return
		}
	}
}

// GetGopsutilMetrics - метод для сбора утилизации ресурсов хоста с определенным интервалом.
func (a *Agent) GetGopsutilMetrics(ctx context.Context, getInterval time.Duration, jobs chan<- GopsutilMetricsSet) {
	defer close(jobs)
	tick := time.NewTicker(getInterval)

	for range tick.C {
		stats, err := mem.VirtualMemory()
		if err != nil {
			logger.Log.Error("failed to collect virtual memory stats", zap.Error(err))
		}
		logger.Log.Info("collect virtual memory stats successfully")
		a.GMetrics.TotalMemory = float64(stats.Total)
		a.GMetrics.FreeMemory = float64(stats.Total)
		a.GMetrics.CPUUtilization = stats.UsedPercent
		select {
		case jobs <- a.GMetrics:
		case <-ctx.Done():
			tick.Stop()
			a.WG.Done()
			return
		}
	}
}

// SendMetrics - метод  через переданный интервал времени передает на сервер метрики.
func (a *Agent) SendMetrics(ctx context.Context, sendInterval time.Duration, jobsMetrics <-chan MetricsSet, jobsGMetrics <-chan GopsutilMetricsSet) {
	tick := time.NewTicker(sendInterval)

	for {
		select {
		case <-tick.C:
			a.Sender.SendToRepo(jobsMetrics, jobsGMetrics)
		case <-ctx.Done():
			tick.Stop()
			a.Sender.SendToRepo(jobsMetrics, jobsGMetrics)
			a.WG.Done()
			return
		}
	}
}
