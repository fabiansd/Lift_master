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

/*func door_timer(timeout chan<- bool, reset <-chan bool, duration time.Duration) {

	var doorOpenTime = duration * time.Second
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-reset:
			timer.Reset(doorOpenTime)
		case <-timer.C:
			timer.Stop()
			timeout <- true
		}
	}
}

/*

func timer_timedOut(timeoutChannel chan bool) bool {
	select {
	case <-timeoutChannel:
		return true
	default:
		return false
	}
}

func main() {
	timeoutChannel := make(chan bool)
	resetChannel := make(chan bool)
	go door_timer(timeoutChannel, resetChannel, 3)
	resetChannel <- true
	for {
		if timer_timedOut(timeoutChannel) {
			fmt.Println("timerout")
			break
		}
	}
}*/
