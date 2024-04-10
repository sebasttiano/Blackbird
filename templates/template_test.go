package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTemplates(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test on panic while parsing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() { ParseTemplates() }, "Parse Templates didn`t panic")
		})
	}
}

func BenchmarkParseTemplates(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ParseTemplates()
	}
}
