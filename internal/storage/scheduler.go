package storage

import (
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"time"
)

func TickerSaver(ticker *time.Ticker, file string) {
	for {
		<-ticker.C
		localStorage := SrvFacility.LocalStorage
		if err := localStorage.SaveToFile(file); err != nil {
			logger.Log.Error("can`t save metrics to file")
		}
		logger.Log.Debug(fmt.Sprintf("save metrics to file %s", file))
	}
}
