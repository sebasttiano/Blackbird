package service

import (
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"time"
)

func TickerSaver(ticker *time.Ticker, service *Service) {
	for {
		<-ticker.C
		if err := service.Save(); err != nil {
			logger.Log.Error("can`t save metrics to file")
		}
		logger.Log.Debug("save metrics to file")
	}
}
