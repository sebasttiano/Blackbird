package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name  string
		level zapcore.Level
	}{
		{name: "INFO", level: zapcore.InfoLevel},
		{name: "DEBUG", level: zapcore.DebugLevel},
		{name: "ERROR", level: zapcore.ErrorLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.level.String())
			if assert.NoError(t, err) {
				assert.Equal(t, Log.Level(), tt.level)
			}
		})
	}
}
