package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"io"
	"math/rand"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

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
	TotalMemory,
	FreeMemory,
	CPUUtilization,
	RandomValue float64
	PollCount int64
}

func main() {
	parseFlags()

	if err := run(); err != nil {
		logger.Log.Error("While executing agent, error occurred", zap.Error(err))
	}
}

func run() error {

	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}
	logger.Log.Info(fmt.Sprintf("Running agent with poll interval %d and report interval %d\n", pollInterval, reportInterval))
	logger.Log.Info(fmt.Sprintf("Metric storage server address is set to %s\n", serverIPAddr))
	mh := NewMetricHandler(pollInterval, reportInterval, "http://"+serverIPAddr, flagSecretKey)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	mh.wg.Add(2)
	go mh.GetMetrics(ctx)
	go mh.GetGopsutilMetrics(ctx)

	g := new(errgroup.Group)

	for i := 0; i < int(flagRateLimit); i++ {
		g.Go(func() error {
			err := mh.IterateStructFieldsAndSend(ctx)
			if err != nil {
				logger.Log.Error("failed to send metrics,", zap.Error(err))
				return err
			}
			return nil
		},
		)
	}
	if err := g.Wait(); err != nil {
		cancel()
	}
	mh.wg.Wait()
	return nil
}

type MetricHandler struct {
	getInterval,
	sendInterval,
	sendCounter time.Duration
	getCounter int64
	client     common.HTTPClient
	rtm        runtime.MemStats
	metrics    MetricsSet
	signKey    string
	wg         sync.WaitGroup
}

func NewMetricHandler(pollInterval, reportInterval int64, serverAddr string, signKey string) MetricHandler {
	getCounter := new(int64)
	return MetricHandler{
		getInterval:  time.Duration(pollInterval) * time.Second,
		sendInterval: time.Duration(reportInterval) * time.Second,
		getCounter:   *getCounter,
		sendCounter:  time.Duration(1) * time.Second,
		client:       common.NewHTTPClient(serverAddr, httpClientRetry, httpClientRetryBackoff),
		signKey:      signKey,
	}
}

// GetMetrics collect runtime metrics
func (m *MetricHandler) GetMetrics(ctx context.Context) {

	tick := time.NewTicker(m.getInterval)

	for {
		select {
		case <-tick.C:
			atomic.AddInt64(&m.getCounter, 1)
			runtime.ReadMemStats(&m.rtm)

			logger.Log.Info("collect memstats successfully")
			m.metrics.Alloc = float64(m.rtm.Alloc)
			m.metrics.TotalAlloc = float64(m.rtm.TotalAlloc)
			m.metrics.BuckHashSys = float64(m.rtm.BuckHashSys)
			m.metrics.Frees = float64(m.rtm.Frees)
			m.metrics.GCCPUFraction = m.rtm.GCCPUFraction
			m.metrics.GCSys = float64(m.rtm.GCSys)
			m.metrics.HeapAlloc = float64(m.rtm.HeapAlloc)
			m.metrics.HeapIdle = float64(m.rtm.HeapIdle)
			m.metrics.HeapInuse = float64(m.rtm.HeapInuse)
			m.metrics.HeapObjects = float64(m.rtm.HeapObjects)
			m.metrics.HeapReleased = float64(m.rtm.HeapReleased)
			m.metrics.HeapSys = float64(m.rtm.HeapSys)
			m.metrics.LastGC = float64(m.rtm.LastGC)
			m.metrics.Lookups = float64(m.rtm.Lookups)
			m.metrics.MCacheInuse = float64(m.rtm.MCacheInuse)
			m.metrics.MCacheSys = float64(m.rtm.MCacheSys)
			m.metrics.MSpanInuse = float64(m.rtm.MSpanInuse)
			m.metrics.MSpanSys = float64(m.rtm.MSpanSys)
			m.metrics.Mallocs = float64(m.rtm.Mallocs)
			m.metrics.NextGC = float64(m.rtm.NextGC)
			m.metrics.NumForcedGC = float64(m.rtm.NumForcedGC)
			m.metrics.NumGC = float64(m.rtm.NumGC)
			m.metrics.OtherSys = float64(m.rtm.OtherSys)
			m.metrics.PauseTotalNs = float64(m.rtm.PauseTotalNs)
			m.metrics.StackInuse = float64(m.rtm.StackInuse)
			m.metrics.StackSys = float64(m.rtm.StackSys)
			m.metrics.Sys = float64(m.rtm.Sys)
			m.metrics.PollCount = m.getCounter
			m.metrics.RandomValue = rand.Float64()
		case <-ctx.Done():
			tick.Stop()
			m.wg.Done()
			return
		}
	}

}

// GetGopsutilMetrics collect gopsutil metrics
func (m *MetricHandler) GetGopsutilMetrics(ctx context.Context) {

	tick := time.NewTicker(m.getInterval)

	for {
		select {
		case <-tick.C:
			stats, err := mem.VirtualMemory()
			if err != nil {
				logger.Log.Error("failed to collect virtual memory stats", zap.Error(err))
			}
			logger.Log.Info("collect virtual memory stats successfully")
			m.metrics.TotalMemory = float64(stats.Total)
			m.metrics.FreeMemory = float64(stats.Total)
			m.metrics.CPUUtilization = stats.UsedPercent
		case <-ctx.Done():
			tick.Stop()
			m.wg.Done()
			return
		}
	}
}

// IterateStructFieldsAndSend prepares url with values and make post request to server
func (m *MetricHandler) IterateStructFieldsAndSend(ctx context.Context) error {

	tick := time.NewTicker(m.sendInterval)

	for {
		select {
		case <-tick.C:
			var metrics models.Metrics
			var metricsBatch []models.Metrics

			value := reflect.ValueOf(m.metrics)
			numFields := value.NumField()
			structType := value.Type()

			for i := 0; i < numFields; i++ {
				field := structType.Field(i)
				fieldValue := value.Field(i)
				metrics.ID = field.Name

				if fieldValue.CanInt() {
					counterVal := fieldValue.Int()
					metrics.Delta = &counterVal
					metrics.MType = "counter"

				} else {
					gaugeVal := fieldValue.Float()
					metrics.Value = &gaugeVal
					metrics.MType = "gauge"
				}

				metricsBatch = append(metricsBatch, metrics)
			}

			if len(metricsBatch) > 0 {
				// Make an HTTP post request
				reqBody, err := json.Marshal(metricsBatch)
				if err != nil {
					logger.Log.Error("couldn`t serialize to json", zap.Error(err))
					return err
				}

				compressedData, err := common.Compress(reqBody)
				if err != nil {
					logger.Log.Error("failed to compress data to gzip", zap.Error(err))
				}

				headers := map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}
				if m.signKey != "" {
					data := *compressedData
					h := hmac.New(sha256.New, []byte(m.signKey))
					if _, err := h.Write(data.Bytes()); err != nil {
						return err
					}
					dst := h.Sum(nil)
					logger.Log.Info("create hmac signature")
					headers["HashSHA256"] = hex.EncodeToString(dst)
				}

				res, err := m.client.Post("/updates/", compressedData, headers)
				if err != nil {
					logger.Log.Error(fmt.Sprintf("couldn`t send metrics batch of length %d", len(metricsBatch)), zap.Error(err))
					continue
				}
				answer, _ := io.ReadAll(res.Body)
				res.Body.Close()

				if res.StatusCode != 200 {
					logger.Log.Error(fmt.Sprintf("error: server return code %d: message: %s", res.StatusCode, answer))
					continue
				}
				logger.Log.Info("send metrics to storage server successfully.")
			}
		case <-ctx.Done():
			tick.Stop()
			return nil
		}
	}
}
