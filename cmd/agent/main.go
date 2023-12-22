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

type Metrics struct {
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

var httpClient HTTPClient

func main() {
	parseFlags()
	httpClient = NewHTTPClient("http://" + serverIPAddr)
	fmt.Printf("Running agent with poll interval %d and report interval %d\n", pollInterval, reportInterval)
	fmt.Printf("Metric storage server address is set to %s\n", serverIPAddr)
	GetMetrics(pollInterval, reportInterval, 30, httpClient)
}

func GetMetrics(pollInterval int64, reportInterval int64, stopLimit int, client HTTPClient) {
	var m Metrics
	var rtm runtime.MemStats
	var getInterval = time.Duration(pollInterval) * time.Second
	var sendInterval = time.Duration(reportInterval) * time.Second
	var getCounter = time.Duration(1) * time.Second
	var sendCounter = time.Duration(1) * time.Second

	m.PollCount = 0

	for i := 0; stopLimit > i; i++ { // TODO make infinite when stoplimit == 0

		time.Sleep(1 * time.Second)

		if getCounter == getInterval {
			runtime.ReadMemStats(&rtm)

			m.Alloc = float64(rtm.Alloc)
			m.TotalAlloc = float64(rtm.TotalAlloc)
			m.BuckHashSys = float64(rtm.BuckHashSys)
			m.Frees = float64(rtm.Frees)
			m.GCCPUFraction = rtm.GCCPUFraction
			m.GCSys = float64(rtm.GCSys)
			m.HeapAlloc = float64(rtm.HeapAlloc)
			m.HeapIdle = float64(rtm.HeapIdle)
			m.HeapInuse = float64(rtm.HeapInuse)
			m.HeapObjects = float64(rtm.HeapObjects)
			m.HeapReleased = float64(rtm.HeapReleased)
			m.HeapSys = float64(rtm.HeapSys)
			m.LastGC = float64(rtm.LastGC)
			m.Lookups = float64(rtm.Lookups)
			m.MCacheInuse = float64(rtm.MCacheInuse)
			m.MCacheSys = float64(rtm.MCacheSys)
			m.MSpanInuse = float64(rtm.MSpanInuse)
			m.MSpanSys = float64(rtm.MSpanSys)
			m.Mallocs = float64(rtm.Mallocs)
			m.NextGC = float64(rtm.NextGC)
			m.NumForcedGC = float64(rtm.NumForcedGC)
			m.NumGC = float64(rtm.NumGC)
			m.OtherSys = float64(rtm.OtherSys)
			m.PauseTotalNs = float64(rtm.PauseTotalNs)
			m.StackInuse = float64(rtm.StackInuse)
			m.StackSys = float64(rtm.StackSys)
			m.Sys = float64(rtm.Sys)
			m.PollCount += 1
			m.RandomValue = rand.Float64()

			getCounter = 0 * time.Second
		}

		if sendCounter == sendInterval {
			if err := iterateStructFieldsAndSend(m, client); err != nil {
				log.Printf("error occured while sending metrics. message: %s", err)
			}
			sendCounter = 0 * time.Second
		}
		getCounter += 1 * time.Second
		sendCounter += 1 * time.Second
	}
}

// iterateStructFieldsAndSend prepares url with values and make post request to server
func iterateStructFieldsAndSend(input interface{}, client HTTPClient) error {

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
		_, err := client.Post(posturl, bytes.NewBuffer([]byte{}), "Content-Type: text/plain")
		if err != nil {
			return err
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
			return nil, errors.New("Check passed header,  it should be in the format '<Name>: <Value>'")
		}

	}
	client := &http.Client{}
	res, err := client.Do(r)

	return res, err
}
