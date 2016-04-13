package queue

import (
	"fmt"
	"operations"
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
}
type ElevatorCost struct {
	Cost int
	Addr string
}

func RecieveCosts(costChan chan operations.Udp_message) {

	OrderCostMap := make(map[Order][]ElevatorCost)

	for {
		select {
		case msg := <-costChan:
			newOrder := Order{msg.Floor, msg.Button}
			newCostReply := ElevatorCost{msg.Cost, msg.Addr}
			/*
				//If
				for oldOrder := range OrderCostMap {
					if oldOrder.Floor == newOrder.Floor && oldOrder.Button == newOrder.Button {
						newOrder = oldOrder
					}
				}*/

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

					OrderCostMap[newOrder] = append(OrderCostMap[newOrder], newCostReply)
				}
			} else { // If not existant, add the new order to the map
				OrderCostMap[newOrder] = []ElevatorCost{newCostReply}
			}

			fmt.Println(OrderCostMap)
			//fmt.Println(len(OrderCostMap))

			/*for order, costReply := range OrderCostMap{
			fmt.Println(order)
			for _,cost := range costReply{
				fmt.Println(cost)
			}*/
		}
	}
}

func AssignOrders(OrderCostMap map[Order][]ElevatorCost) {

}

func RemoteAssignOrder() {

}
