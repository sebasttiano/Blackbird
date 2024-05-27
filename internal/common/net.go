package common

import (
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	"net"
)

func GetLocalIP(socket string) (net.IP, error) {
	conn, err := net.Dial("udp", socket)
	if err != nil {
		logger.Log.Error("failed to lookup local addr", zap.Error(err))
		return nil, nil
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP, nil
}
