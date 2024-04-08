package repository

import (
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"time"
)

func TickerSaver(ticker *time.Ticker, store Store) {
	for {
		<-ticker.C
		if err := store.Save(); err != nil {
			logger.Log.Error("can`t save metrics to file")
		}
		logger.Log.Debug("save metrics to file")
	}
}
