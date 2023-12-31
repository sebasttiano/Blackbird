package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
