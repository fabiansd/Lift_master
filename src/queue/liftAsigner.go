package queue

import (
	"fmt"
)

type Direction int 
const (
	DIRN_DOWN = -1 + iota
	DIRN_STOP
	DIRN_UP
)

type Elevator struct {
	Floor     int
	PrevFloor int
	Dir       Direction
}

type Udp_message struct {
	Category int
	Floor    int
	Button   int
	Cost     int
	Addr     string `json:"-"`
}

type Order struct {
	Floor int
	Button int
}
type ElevatorCost struct {
	Cost int
	Addr string
}

var msg1 = Udp_message{Category:1,Floor:1,Button:0,Cost:10,Addr:"lift1"}
var msg11 = Udp_message{Category:1,Floor:1,Button:0,Cost:9,Addr:"lift1"}
var msg2 = Udp_message{Category:1,Floor:1,Button:0,Cost:11,Addr:"lift1"}
var msg3 = Udp_message{Category:1,Floor:4,Button:0,Cost:12,Addr:"lift3"}
var msg4 = Udp_message{Category:1,Floor:2,Button:0,Cost:12,Addr:"lift4"}
var msg5 = Udp_message{Category:1,Floor:1,Button:1,Cost:10,Addr:"lift3"}
var msg6 = Udp_message{Category:1,Floor:4,Button:0,Cost:10,Addr:"lift3"}
var msg7 = Udp_message{Category:1,Floor:4,Button:0,Cost:10,Addr:"lift3"}


var liftsOnline int = 3

func main() {
	
	var in = make(chan Udp_message)

	go liftassigner(in)
	in <- msg1
	in <- msg11
	in <- msg2
	in <- msg3
	in <- msg4
	in <- msg5
	in <- msg6
	in <- msg7
}

func liftassigner(in chan Udp_message){

	OrderCostMap := make(map[Order][]ElevatorCost)


	for{
		select{
			case msg := <- in :
				newOrder := Order{msg.Floor,msg.Button}
				newCostReply := ElevatorCost{msg.Cost,msg.Addr}
/*
				//If 
				for oldOrder := range OrderCostMap {
					if oldOrder.Floor == newOrder.Floor && oldOrder.Button == newOrder.Button {
						newOrder = oldOrder
					}
				}*/

				// check if the order order exists in the map
				if CostReply,err := OrderCostMap[newOrder];err{
					//fmt.Println(CostReply)
					exists := false
					//check if ordercost is already registred
					for n,cost := range CostReply{
						if cost.Cost != newCostReply.Cost && cost.Addr == newCostReply.Addr{
							//Deletes the cost element in the order-costlist to replace it if the cost 
							//has been updatet for a particular elavtor at this particular order
							OrderCostMap[newOrder] = append(OrderCostMap[newOrder][:n],OrderCostMap[newOrder][n+1:]...)
						}else if cost == newCostReply{
							exists = true
						}
					}
					//Register order if the cost is unregistred
					if !exists{

						OrderCostMap[newOrder] = append(OrderCostMap[newOrder],newCostReply)
					}
				}else{ // If not existant, add the new order to the map
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
