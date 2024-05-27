package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{name: "Check localhost", addr: "127.0.0.1:80", want: "127.0.0.1"},
		{name: "Check google dns", addr: "8.8.8.8:80", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localAddr, err := GetLocalIP(tt.addr)
			if assert.NoError(t, err) {
				if tt.want != "" {
					assert.Equal(t, tt.want, localAddr.String())
				}
			}
		})
	}
}
