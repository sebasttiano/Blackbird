package main

import (
	"testing"

	"github.com/sebasttiano/Blackbird.git/internal/service"
)

func BenchmarkApp_Initialize(b *testing.B) {
	b.ReportAllocs()
	app := newApp()
	for i := 0; i < b.N; i++ {
		app.Initialize(&service.ServiceSettings{DBSave: false}, "SECRET_KEY")
	}
}
