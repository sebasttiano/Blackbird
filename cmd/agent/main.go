package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

var pollInterval int = 2
var reportInterval int64 = 10

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

func main() {
	GetMetrics(pollInterval)
}

func GetMetrics(pollInterval int) {
	var m Metrics
	var rtm runtime.MemStats
	var interval = time.Duration(pollInterval) * time.Second
	var startTime = time.Now().Unix()
	m.PollCount = 0

	for {
		time.Sleep(interval)
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

		if (startTime-time.Now().Unix())%reportInterval == 0 {
			iterateStructFieldsAndSend(m)
		}
	}
}

// iterateStructFieldsAndSend prepares url with values and make post request to server
func iterateStructFieldsAndSend(input interface{}) {

	fmt.Println(time.Now().Unix())
	var posturl string

	value := reflect.ValueOf(input)
	numFields := value.NumField()
	structType := value.Type()

	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)
		if field.Name == "PollCount" {
			posturl = fmt.Sprintf("http://localhost:8080/update/counter/%s/%d", field.Name, fieldValue)
		} else {
			posturl = fmt.Sprintf("http://localhost:8080/update/gauge/%s/%0.f", field.Name, fieldValue)
		}

		// Create an HTTP post request
		r, err := http.NewRequest("POST", posturl, bytes.NewBuffer([]byte{}))
		if err != nil {
			panic(err)
		}
		r.Header.Add("Content-Type", "text/plain")
		client := &http.Client{}
		res, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
	}
}
