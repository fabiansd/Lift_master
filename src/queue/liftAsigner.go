package queue

import (
	//"driver"
	"fmt"
	"operations"
	"time"
)

/*
func RecieveCosts(costChan chan operations.Udp_message) {
	for {
		msg := <-costChan
		fmt.Println("Cost")
		fmt.Println(msg.Cost)
		operations.Fsm_neworder(int(msg.Floor), msg.Button)
	}
}*/

type Order struct {
	Floor  int
	Button int
	timer  *time.Timer
}
type ElevatorCost struct {
	Cost int
	Addr string
}

var elevatorsOnline int

func RecieveCosts(costChan chan operations.Udp_message, newOrderChan chan operations.Keypress, orderCompleteChannel chan Order, messageOut chan operations.Udp_message, elevatorsOnlineChan chan int) {
	waitForOrderChan := make(chan Order, 10)
	OrderCostMap := make(map[Order][]ElevatorCost)
	newOrderTimeout := make(chan *Order)

	//const timeoutConst = 3*time.Second
	for {
		select {
		case msg := <-costChan:
			newOrder := Order{Floor: msg.Floor, Button: msg.Button}
			newCostReply := ElevatorCost{msg.Cost, msg.Addr}

			//If
			for oldOrder := range OrderCostMap {
				if oldOrder.Floor == newOrder.Floor && oldOrder.Button == newOrder.Button {
					newOrder = oldOrder
				}
			}

			// check if the  order exists in the map
			if CostReply, err := OrderCostMap[newOrder]; err {
				//fmt.Println(CostReply)
				exists := false
				//check if ordercost is already registred
				for n, cost := range CostReply {
					if cost.Cost != newCostReply.Cost && cost.Addr == newCostReply.Addr {
						//Deletes the cost element in the order-costlist to replace it if the cost
						//has been updatet for a particular elavtor at this particular order
						OrderCostMap[newOrder] = append(OrderCostMap[newOrder][:n], OrderCostMap[newOrder][n+1:]...)
					} else if cost == newCostReply {
						exists = true
					}
				}
				//Register order if the cost is unregistred
				if !exists {
					newOrder.timer.Reset(4 * time.Second)
					OrderCostMap[newOrder] = append(OrderCostMap[newOrder], newCostReply)

				}
			} else { // If not existant, add the new order to the map

				newOrder.timer = time.NewTimer(4 * time.Second)
				go costreplyTimer(newOrderTimeout, &newOrder)
				OrderCostMap[newOrder] = []ElevatorCost{newCostReply}
			}
			fmt.Println("MAP: ", OrderCostMap)
			AssignOrders(OrderCostMap, newOrder, false, messageOut, waitForOrderChan)

		case no := <-newOrderTimeout:

			AssignOrders(OrderCostMap, *no, true, messageOut, waitForOrderChan)

			//case waitForOrderCompleted := <-waitForOrderChan: //We know the ordercomplete-message reaches this stage
			//go waitForCompletion(waitForOrderCompleted, newOrderChan, orderCompleteChannel)
		case eo := <-elevatorsOnlineChan:
			elevatorsOnline = eo
			fmt.Println("the update found its way to liftassigner: ", elevatorsOnline)
		}

	}
}

/*
func waitForCompletion(order Order, newOrderChan chan operations.Keypress, orderCompleteChannel chan Order) {

	waitTimer := time.NewTimer(time.Second * 5)
	fmt.Println("wait for ordercomplete-timer started")

	select {
	case OC := <-orderCompleteChannel:
		if order == OC {
			fmt.Println("ordercompleted!", OC, "ORDER KILLED")

		}
	case <-waitTimer.C:
		newOrderChan <- operations.Keypress{Floor: order.Floor, Button: order.Button}
		fmt.Println("Resent to neworderchannel for re-broadcast")

	}
}*/

func AssignOrders(OrderCostMap map[Order][]ElevatorCost, no Order, isCostreplyTimeout bool, messageOut chan operations.Udp_message, waitForOrderChan chan Order) {
	//fmt.Println("no from timeout", no)
	//Må kont oppdateres
	for order, costList := range OrderCostMap { //for each costlist

		if (isCostreplyTimeout && order == no) || len(costList) == elevatorsOnline {

			smallestCost := 10000
			scostElevator := ""
			for _, cost := range costList {
				if smallestCost > cost.Cost { //Assign to the smallest cost
					smallestCost = cost.Cost
					scostElevator = cost.Addr
				} else if smallestCost == cost.Cost { //Assign to the biggest IP
					if scostElevator > cost.Addr {
						smallestCost = cost.Cost
						scostElevator = cost.Addr
					}
				}
			}

			order.timer.Stop()
			//waitForOrderChan <- order
			fmt.Println("Order:{", order.Floor, " ", order.Button, "with cost: ", smallestCost, "} assigned to", scostElevator)
			backupRequest[order.Floor][order.Button] = true
			fmt.Println("The backuplist: ", backupRequest)
			messageOut <- operations.Udp_message{Category: operations.AssignedOrder, Floor: order.Floor, Button: order.Button, Cost: smallestCost, AssignAddr: scostElevator}
			delete(OrderCostMap, order)
		}
	}
}

func costreplyTimer(newOrderTimeout chan<- *Order, newOrder *Order) {
	fmt.Println("costReply timeout timer started")
	<-newOrder.timer.C
	newOrderTimeout <- newOrder
}
