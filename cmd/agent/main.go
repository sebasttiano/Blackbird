package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strings"
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
	fmt.Printf("Running agent with poll interval %d and report interval %d\n", pollInterval, reportInterval)
	fmt.Printf("Metric storage server address is set to %s\n", serverIPAddr)
	mh := NewMetricHandler(pollInterval, reportInterval, 30, "http://"+serverIPAddr)
	if err := mh.GetMetrics(); err != nil {
		log.Println("error in getmetrics occured")
	}
}

type MetricHandler struct {
	getInterval,
	sendInterval,
	getCounter,
	sendCounter time.Duration
	stopLimit int
	client    HTTPClient
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
		client:       NewHTTPClient(serverAddr),
	}
}

func (m *MetricHandler) GetMetrics() error {

	var err error
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
			if err = IterateStructFieldsAndSend(m.metrics, m.client); err != nil {
				return err
			}
			m.sendCounter = 0 * time.Second
		}

		m.getCounter += 1 * time.Second
		m.sendCounter += 1 * time.Second
	}
	return nil
}

// IterateStructFieldsAndSend prepares url with values and make post request to server
func IterateStructFieldsAndSend(input interface{}, client HTTPClient) error {

	var posturl string

	value := reflect.ValueOf(input)
	numFields := value.NumField()
	structType := value.Type()

	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)
		if fieldValue.CanInt() {
			posturl = fmt.Sprintf("/update/counter/%s/%d", field.Name, fieldValue.Int())

		} else {
			posturl = fmt.Sprintf("/update/gauge/%s/%0.f", field.Name, fieldValue.Float())

		}

		// Make an HTTP post request
		res, err := client.Post(posturl, bytes.NewBuffer([]byte{}), "Content-Type: text/plain")
		if err != nil {
			return err
		}
		res.Body.Close()

		if res.StatusCode != 200 {
			return fmt.Errorf("error: server return code %d, while sending metric %s", res.StatusCode, field.Name)
		}
	}
	return nil
}

// HTTPClient simple client
type HTTPClient struct {
	url string
}

func NewHTTPClient(url string) HTTPClient {
	return HTTPClient{url: url}
}

// Post implements http post requests
func (c HTTPClient) Post(urlSuffix string, body io.Reader, header string) (*http.Response, error) {

	r, err := http.NewRequest("POST", c.url+urlSuffix, body)
	if err != nil {
		return nil, err
	}
	if header != "" {
		splitHeader := strings.Split(header, ":")
		if len(splitHeader) == 2 {
			r.Header.Add(splitHeader[0], splitHeader[1])
		} else {
			return nil, errors.New("error: check passed header,  it should be in the format '<Name>: <Value>'")
		}

	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}
