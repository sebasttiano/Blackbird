package main

import (
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"testing"
)

func BenchmarkApp_Initialize(b *testing.B) {
	app := newApp()
	for i := 0; i < b.N; i++ {
		app.Initialize(&service.ServiceSettings{DBSave: false}, "SECRET_KEY")
	}
}
