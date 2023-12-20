package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type AgentArgs struct {
	serverIpAddr   string
	pollInterval   int64
	reportInterval int64
}

func Test_parseFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want AgentArgs
	}{
		{name: "Check parsing variables",
			args: []string{"myFile", "-a", "localhost:9000", "-p", "4", "-r", "20"},
			want: AgentArgs{serverIpAddr: "localhost:9000", pollInterval: 4, reportInterval: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			parseFlags()
			assert.Equal(t, tt.want.serverIpAddr, serverIPAddr)
			assert.Equal(t, tt.want.pollInterval, pollInterval)
			assert.Equal(t, tt.want.reportInterval, reportInterval)
		})
	}
}
