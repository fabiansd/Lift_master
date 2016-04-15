package queue

import (
	"fmt"
	"time"
	"operations"
)
//Nye channels
var backupRequest [4][3]string //Alle bestillinger som blir assigned blir registrert her som true
var backupChan = make(chan Order,100)
var timeoutChan = make(chan Order,100)

func BackupHandler(newOrderChan chan operations.Keypress, orderCompleteChannel chan Order, outgoingMsg chan operations.Udp_message) {

	for {
		select{
		case nO := <- backupChan:
				if (backupRequest[nO.Floor][nO.Button] == ""){// "" = false
				backupRequest[nO.Floor][nO.Button] = nO.Addr
				fmt.Println("Backup updated: ", backupRequest)
				go completionTimer(nO)
			}

			
		case timeout := <- timeoutChan:
			if backupRequest[timeout.Floor][timeout.Button] != ""{
				newOrderChan <- operations.Keypress{timeout.Floor,timeout.Button}
				//outgoingMsg <- operations.Udp_message{Category: operations.Killfeed,AssignAddr:timeout.Addr}
				backupRequest[timeout.Floor][timeout.Button] = ""
				fmt.Println("Resend the order ", timeout, " to ", timeout.Addr)

			}
		case ordeComp := <- orderCompleteChannel:
			if backupRequest[ordeComp.Floor][ordeComp.Button] != ""{
				backupRequest[ordeComp.Floor][ordeComp.Button] = ""
				fmt.Println("Order deleted: ", ordeComp, " Backup updated: ", backupRequest)
			}
		default:
		}

	}
}


func completionTimer(order Order){
	fmt.Println("completionTimer started!")
	time.Sleep(15 * time.Second)
	timeoutChan <- order
}
