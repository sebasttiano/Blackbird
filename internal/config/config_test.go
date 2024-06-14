package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
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

func Test_parseServerFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Config
	}{
		{name: "Check parsing variables",
			args: []string{"myServer", "-f", "/files/myfile.txt", "-i", "4", "-k", "secret", "-crypto-key", "/tmp/key", "-g", ":3200"},
			want: Config{FileStoragePath: "/files/myfile.txt", StoreInterval: 4, SecretKey: "secret", CryptoKey: "/tmp/key", GRPSServerIPAddr: ":3200"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			//config := parseServerFlags()  TODO: test parsing servers also
			//assert.Equal(t, config, tt.want)

		})
	}
}

func TestNewServerConfig(t *testing.T) {
	f, _ := strconv.ParseBool("false")
	y, _ := strconv.ParseBool("true")
	test := struct {
		name string
		want *Config
	}{
		name: "default", want: &Config{
			ServerIPAddr:    "localhost:8080",
			PollInterval:    2,
			ReportInterval:  5,
			RateLimit:       1,
			Profiler:        &f,
			StoreInterval:   300,
			RestoreMetrics:  &y,
			FileStoragePath: "/tmp/metrics-db.json",
		},
	}
	t.Run(test.name, func(t *testing.T) {
		config := &Config{}
		config.SetDefault()
		assert.Equalf(t, test.want, config, "NewServerConfig()")
	})
}
