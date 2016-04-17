package globalOperations

import (
	"elevatorOperations"
	"fmt"
	"time"
)

var backupRequest [4][3]bool 
var backupChan = make(chan Order, 100)
var timeoutChan = make(chan Order, 100)

func BackupHandler(newOrderChan chan elevatorOperations.Keypress, orderCompleteChannel chan Order, outgoingMsg chan elevatorOperations.Udp_message) {

	for {
		select {
		case nO := <-backupChan:
			if backupRequest[nO.Floor][nO.Button] == false { 
				backupRequest[nO.Floor][nO.Button] = true
				elevatorOperations.Fsm_printrequest(backupRequest, "The backup requests")
				go completionTimer(nO)
			}

		case timeout := <-timeoutChan:
			if backupRequest[timeout.Floor][timeout.Button] == true {
				newOrderChan <- elevatorOperations.Keypress{timeout.Floor, timeout.Button}
				backupRequest[timeout.Floor][timeout.Button] = false
				fmt.Println(elevatorOperations.Yellow, "Resent the order ", timeout, " to ", timeout.Addr, elevatorOperations.White)

			}
		case ordeComp := <-orderCompleteChannel:
			if backupRequest[ordeComp.Floor][ordeComp.Button] == true {
				backupRequest[ordeComp.Floor][ordeComp.Button] = false
				elevatorOperations.Fsm_printrequest(backupRequest, "The backup requests")
			}
		default:
		}

	}
}

func completionTimer(order Order) {
	time.Sleep(15 * time.Second)
	timeoutChan <- order
}
