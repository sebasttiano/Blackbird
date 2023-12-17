package main

import (
	"fmt"
	"runtime"
)

type Metrics struct {
	Alloc,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle float64
}

func main() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	fmt.Println(rtm)
}

//func NewMetrics(d)
