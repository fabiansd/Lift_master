package queue

import (
	//"fmt"
	"driver"
)

var backupRequest [driver.N_FLOORS][driver.N_BUTTONS - 1]bool //baskup matrix for global orders (up and down)
