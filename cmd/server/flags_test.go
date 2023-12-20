package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type ServerArgs struct {
	flagRunAddr string
}

func Test_parseFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want ServerArgs
	}{
		{name: "Check parsing variable",
			args: []string{"myFile", "-a", "localhost:9001"},
			want: ServerArgs{flagRunAddr: "localhost:9001"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			parseFlags()
			assert.Equal(t, tt.want.flagRunAddr, flagRunAddr)
		})
	}
}
