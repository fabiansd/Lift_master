package operations

import (
	"time"
)

var endTime time.Time
var timerActive bool

func Timer_start() {
	endTime = time.Now().Add(3 * time.Second)
	timerActive = true
}

func Timer_stop() {
	timerActive = false
}

func Timer_timedout() bool {
	return (timerActive && time.Now().After(endTime))
}
