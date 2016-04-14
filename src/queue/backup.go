package main

import (
	"fmt"
	"time"
)
//Nye channels
var backupRequest [4][3]bool //Alle bestillinger som blir assigned blir registrert her som true
var backupChan = make(chan Order,100)
var timeoutChan = make(chan Order,100)

//Fra main: (MÅ GJØRE NEWORDERCHAN GLOBAL I MAIN)
var newOrderChan = make(chan Order,100)
var orderCompletedChan = make(chan Order,100)

//Fra liftassigner
type Order struct {
	Floor int
	Button int
}


func main() {

	go sendInBackups()
	go readNeworderChan()
	for {
		select{
		case nO := <- backupChan:
				if backupRequest[nO.Floor][nO.Button] == false{
				backupRequest[nO.Floor][nO.Button] = true
			}
			fmt.Println("Backup recieved: ", nO)
			go completionTimer(nO)
		case timeout := <- timeoutChan:
			if backupRequest[timeout.Floor][timeout.Button] == true{
				fmt.Println("Resend the order ", timeout, " to newOrderChan")
				newOrderChan <- Order{timeout.Floor,timeout.Button}
			}else if backupRequest[timeout.Floor][timeout.Button] == false{
				fmt.Println("The order is completed, no worries ;)")
			}
		case ordeComp := <- orderCompletedChan:
			if backupRequest[ordeComp.Floor][ordeComp.Button] == true{
				backupRequest[ordeComp.Floor][ordeComp.Button] = false
				fmt.Println("Order deleted: ", ordeComp, " From: ", backupRequest)
			}
		default:
		}
	}
}


func completionTimer(order Order){
	fmt.Println("completionTimer stared")
	time.Sleep(5 * time.Second)
	timeoutChan <- order
}

func readNeworderChan(){
	for{
		a := <- newOrderChan
		fmt.Println("NEW! Floor: ", a.Floor, " Button: ", a.Button)
		//backupChan <- Order{Floor: a.Floor,Button: a.Button}
	}
}

func sendInBackups(){
	time.Sleep(1 * time.Second)
	backupChan <- Order{Floor: 1,Button: 1}
	time.Sleep(3 * time.Second)
	orderCompletedChan <- Order{Floor: 1,Button: 1}


}
