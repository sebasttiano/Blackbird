package agent

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
	"io"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
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
	RandomValue float64
	PollCount int64
}

type GopsutilMetricsSet struct {
	TotalMemory,
	FreeMemory,
	CPUUtilization float64
}

type Agent struct {
	getCounter int64
	client     common.HTTPClient
	rtm        runtime.MemStats
	signKey    string
	Metrics    MetricsSet
	GMetrics   GopsutilMetricsSet
	WG         sync.WaitGroup
}

func NewAgent(serverAddr string, clientRetries int, backoffFactor uint, signKey string) *Agent {
	getCounter := new(int64)
	return &Agent{
		getCounter: *getCounter,
		client:     common.NewHTTPClient(serverAddr, clientRetries, backoffFactor),
		signKey:    signKey,
	}
}

// GetMetrics collect runtime metrics
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

// GetGopsutilMetrics collect gopsutil metrics
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

// IterateStructFieldsAndSend prepares url with values and make post request to server
func (a *Agent) IterateStructFieldsAndSend(ctx context.Context, sendInterval time.Duration, jobsMetrics <-chan MetricsSet, jobsGMetrics <-chan GopsutilMetricsSet) {

	tick := time.NewTicker(sendInterval)

	for {
		select {
		case <-tick.C:
			var metric MetricsSet
			var metricG GopsutilMetricsSet
			var metrics models.Metrics
			var metricsBatch []models.Metrics
			var value reflect.Value

			select {
			case metric = <-jobsMetrics:
				value = reflect.ValueOf(metric)
			case metricG = <-jobsGMetrics:
				value = reflect.ValueOf(metricG)
			}
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
					continue
				}

				compressedData, err := common.Compress(reqBody)
				if err != nil {
					logger.Log.Error("failed to compress data to gzip", zap.Error(err))
					continue
				}

				headers := map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}
				if a.signKey != "" {
					data := *compressedData
					h := hmac.New(sha256.New, []byte(a.signKey))
					if _, err := h.Write(data.Bytes()); err != nil {
						logger.Log.Error("failed to create hmac signature")
						continue
					}
					dst := h.Sum(nil)
					logger.Log.Info("create hmac signature")
					headers["HashSHA256"] = hex.EncodeToString(dst)
				}

				res, err := a.client.Post("/updates/", compressedData, headers)
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
			a.WG.Done()
			return
		}
	}
}
