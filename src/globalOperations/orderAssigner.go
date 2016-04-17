package globalOperations

import (
	"elevatorOperations"
	"time"
)

type Order struct {
	Floor  int
	Button int
	timer  *time.Timer
	Addr   string
}
type ElevatorCost struct {
	Cost int
	Addr string
}

var elevatorsOnline int

func RecieveCosts(costChan chan elevatorOperations.Udp_message, messageOut chan elevatorOperations.Udp_message, elevatorsOnlineChan chan int) {
	waitForOrderChan := make(chan Order, 10)
	OrderCostMap := make(map[Order][]ElevatorCost)
	newOrderTimeout := make(chan *Order)

	for {
		select {
		case msg := <-costChan:
			newOrder := Order{Floor: msg.Floor, Button: msg.Button}
			newCostReply := ElevatorCost{msg.Cost, msg.Addr}

			for oldOrder := range OrderCostMap {
				if oldOrder.Floor == newOrder.Floor && oldOrder.Button == newOrder.Button {
					newOrder = oldOrder
				}
			}

			if CostReply, err := OrderCostMap[newOrder]; err {
				exists := false
				for n, cost := range CostReply {
					if cost.Cost != newCostReply.Cost && cost.Addr == newCostReply.Addr {
						//Updates cost if elevator cost for an order is updated
						OrderCostMap[newOrder] = append(OrderCostMap[newOrder][:n], OrderCostMap[newOrder][n+1:]...)
					} else if cost == newCostReply {
						exists = true
					}
				}
				if !exists {
					newOrder.timer.Reset(4 * time.Second)
					OrderCostMap[newOrder] = append(OrderCostMap[newOrder], newCostReply)

				}
			} else { // If non-existent, add the new order to the map

				newOrder.timer = time.NewTimer(4 * time.Second)
				go costreplyTimer(newOrderTimeout, &newOrder)
				OrderCostMap[newOrder] = []ElevatorCost{newCostReply}
			}
			AssignOrders(OrderCostMap, newOrder, false, messageOut, waitForOrderChan)
		case no := <-newOrderTimeout:
			AssignOrders(OrderCostMap, *no, true, messageOut, waitForOrderChan)
		case eo := <-elevatorsOnlineChan:
			elevatorsOnline = eo
		}

	}
}

func AssignOrders(OrderCostMap map[Order][]ElevatorCost, no Order, isCostreplyTimeout bool, messageOut chan elevatorOperations.Udp_message, waitForOrderChan chan Order) {
	for order, costList := range OrderCostMap {

		if (isCostreplyTimeout && order == no) || len(costList) == elevatorsOnline {
			
			//Chooses the lowest cost, or the elevator with highest IP, in the ordercostmap 
			smallestCost := 10000
			scostElevator := ""
			for _, cost := range costList {
				if smallestCost > cost.Cost {
					smallestCost = cost.Cost
					scostElevator = cost.Addr
				} else if smallestCost == cost.Cost {
					if scostElevator > cost.Addr {
						smallestCost = cost.Cost
						scostElevator = cost.Addr
					}
				}
			}

			order.timer.Stop()
			backupChan <- Order{Floor: order.Floor, Button: order.Button, Addr: scostElevator}
			messageOut <- elevatorOperations.Udp_message{Category: elevatorOperations.AssignedOrder, Floor: order.Floor, Button: order.Button, Cost: smallestCost, AssignAddr: scostElevator}
			delete(OrderCostMap, order)
		}
	}
}

func costreplyTimer(newOrderTimeout chan<- *Order, newOrder *Order) {
	<-newOrder.timer.C
	newOrderTimeout <- newOrder
}
