package common

import (
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"time"
)

func Schedule(ticker *time.Ticker, file string) {
	for {
		<-ticker.C

		localStorage := *storage.GetCurrentStorage()
		if err := localStorage.SaveToFile(file); err != nil {
			logger.Log.Error("can`t save metrics to file")
		}
		logger.Log.Debug(fmt.Sprintf("save metrics to file %s", file))
	}
}
