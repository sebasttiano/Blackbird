package common

import (
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"time"
)

func Schedule(ticker *time.Ticker) {
	for {
		// This blocks until a value is received, the ticker
		// sends a value to it every one minute (or the interval specified)
		<-ticker.C

		data := storage.GetCurrentStorage()
		fmt.Println(*data)
		fmt.Println("Tick")
	}
}
