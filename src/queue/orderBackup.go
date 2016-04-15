package queue

import (
	"fmt"
	"time"
	"operations"
)
//Nye channels
var backupRequest [4][3]bool //Alle bestillinger som blir assigned blir registrert her som true
var backupChan = make(chan Order,100)
var timeoutChan = make(chan Order,100)

func BackupHandler(newOrderChan chan operations.Keypress, orderCompleteChannel chan Order, outgoingMsg chan operations.Udp_message) {

	for {
		select{
		case nO := <- backupChan:
				if (backupRequest[nO.Floor][nO.Button] == false){// "" = false
				backupRequest[nO.Floor][nO.Button] = true
				operations.Fsm_printrequest(backupRequest, "The backup requests")
				go completionTimer(nO)
			}

			
		case timeout := <- timeoutChan:
			if backupRequest[timeout.Floor][timeout.Button] == true{
				newOrderChan <- operations.Keypress{timeout.Floor,timeout.Button}
				backupRequest[timeout.Floor][timeout.Button] = false
				fmt.Println(operations.Yellow ,"Resent the order ", timeout, " to ", timeout.Addr, operations.White)

			}
		case ordeComp := <- orderCompleteChannel:
			if backupRequest[ordeComp.Floor][ordeComp.Button] == true{
				backupRequest[ordeComp.Floor][ordeComp.Button] = false
				operations.Fsm_printrequest(backupRequest, "The backup requests")
			}
		default:
		}

	}
}


func completionTimer(order Order){
	time.Sleep(15 * time.Second)
	timeoutChan <- order
}
