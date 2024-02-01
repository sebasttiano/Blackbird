package main

import (
	"encoding/json"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"go.uber.org/zap"
	"math/rand"
	"reflect"
	"runtime"
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
	mh := NewMetricHandler(pollInterval, reportInterval, 600, "http://"+serverIPAddr)
	if err := mh.GetMetrics(); err != nil {
		logger.Log.Error("error in getmetrics occured", zap.Error(err))
	}
	return nil
}

type MetricHandler struct {
	getInterval,
	sendInterval,
	getCounter,
	sendCounter time.Duration
	stopLimit int
	client    common.HTTPClient
	rtm       runtime.MemStats
	metrics   MetricsSet
}

func NewMetricHandler(pollInterval, reportInterval int64, stopLimit int, serverAddr string) MetricHandler {
	return MetricHandler{
		getInterval:  time.Duration(pollInterval) * time.Second,
		sendInterval: time.Duration(reportInterval) * time.Second,
		getCounter:   time.Duration(1) * time.Second,
		sendCounter:  time.Duration(1) * time.Second,
		stopLimit:    stopLimit,
		client:       common.NewHTTPClient(serverAddr, httpClientRetryTimeout, httpClientRetry),
	}
}

func (m *MetricHandler) GetMetrics() error {

	m.metrics.PollCount = 0

	for i := 0; m.stopLimit > i; i++ { // TODO make infinite when stoplimit == 0

		time.Sleep(1 * time.Second)

		if m.getCounter == m.getInterval {
			runtime.ReadMemStats(&m.rtm)

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
			m.metrics.PollCount += 1
			m.metrics.RandomValue = rand.Float64()

			m.getCounter = 0 * time.Second
		}

		if m.sendCounter == m.sendInterval {
			if err := IterateStructFieldsAndSend(m.metrics, m.client); err != nil {
				logger.Log.Error("failed to send metrics to server. error:", zap.Error(err))
				continue
			}
			m.sendCounter = 0 * time.Second
		}

		m.getCounter += 1 * time.Second
		m.sendCounter += 1 * time.Second
	}
	return nil
}

// IterateStructFieldsAndSend prepares url with values and make post request to server
func IterateStructFieldsAndSend(input interface{}, client common.HTTPClient) error {

	var metrics models.Metrics
	var metricsBatch []models.Metrics

	value := reflect.ValueOf(input)
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

		res, err := client.Post("/updates/", compressedData, map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"})

		if err != nil {
			logger.Log.Error(fmt.Sprintf("couldn`t send metrics batch of length %d", len(metricsBatch)), zap.Error(err))
			return err
		}
		res.Body.Close()

		if res.StatusCode != 200 {
			return fmt.Errorf("error: server return code %d, while sending metric batch", res.StatusCode)
		}
	}
	return nil
}
