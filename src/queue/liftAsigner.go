package queue

import (
	"fmt"
	"operations"
)

func RecieveCosts(costChan chan operations.Udp_message) {
	for {
		msg := <-costChan
		fmt.Println("Cost")
		fmt.Println(msg.Cost)
		operations.Fsm_neworder(int(msg.Floor), msg.Button)
	}
}
