package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type AgentArgs struct {
	flagRunAddr    string
	pollInterval   int64
	reportInterval int64
	flagRateLimit  uint64
	flagSecretKey  string
}

func Test_parseAgentFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want AgentArgs
	}{
		{name: "Check parsing variables",
			args: []string{"myFile", "-a", "localhost:9001", "-p", "4", "-r", "20", "-l", "3", "-k", "secret-key"},
			want: AgentArgs{flagRunAddr: "localhost:9001", pollInterval: 4, reportInterval: 20, flagRateLimit: 3, flagSecretKey: "secret-key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			config := parseAgentFlags()
			assert.Equal(t, tt.want.flagRunAddr, config.ServerIPAddr)
			assert.Equal(t, tt.want.pollInterval, config.PollInterval)
			assert.Equal(t, tt.want.reportInterval, config.ReportInterval)
			assert.Equal(t, tt.want.flagRateLimit, config.RateLimit)
			assert.Equal(t, tt.want.flagSecretKey, config.SecretKey)
		})
	}
}
