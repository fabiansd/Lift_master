package main



import (
	"fmt"
	//"math"
	)

const (
	N_FLOORS = 4;
	travelTime = 1;
	stopTime = 3;
)

var (
	//cart1{dir = 1, pos = 0} Cart;
	//cart2{dir = 0, pos = 1;} Cart;
	orders Orders;
	
)

type Orders struct{
	orders[4][2] int
}

func (self *Orders) addOrder(floor int, dir int) {
	if dir == -1{
		dir = 0;
	}
	self.orders[floor][dir] = 1;
}

func (self *Orders) removeOrder(floor int, dir int) {
	self.orders[floor][dir] = 0;
}

func (self Orders) checkIfOrder(floor int, dir int) int{
	if dir == -1 {
		dir = 0;
	}
	if self.orders[floor][dir] == 1{
		return 1;
	}else{
		return 0;
	}
}

type Cart struct{
	pos int
	dir int 
	commands[4] int
	idnumber int
}

func (self *Cart) addCommand(floor int) {
	self.commands[floor] = 1;
}

func (self *Cart) removeCommand(floor int) {
	self.commands[floor] = 0;
}

func (self Cart) checkIfCommand(floor int) int {
	if self.commands[floor] == 1{
		return 1;
	}else{
		return 0;
	}
}


func (self Cart) curDir() int {
	return self.dir;
}

func (self Cart) curPos() int {
	return self.pos;
}

func (self Cart) id() int {
	return self.idnumber;
}

func (self Cart) checkIfStopOnFloor(orderFloor int, orderDirection int, orders Orders) int{
	if self.checkIfCommand(orderFloor) == 1{
		return 1;
	}else if orders.checkIfOrder(orderFloor, orderDirection) == 1{
		return 1;
	}else{
		return 0;
	}
}

func abs(x int) int{
	if x > 0{
		return x;
	}else if x < 0{
		return -x;
	}else{
		return 0
	}
}

func addedCostForElevator(orders Orders, cart Cart, orderFloor int, orderDirection int) int{
	turnFloor := cart.curPos();
	lastFloor := turnFloor;
	if cart.curDir() == 0{
		return abs(cart.curPos() - orderFloor) * travelTime
	}
	for floor := cart.curPos(); floor < N_FLOORS; floor += cart.curDir() {
		if floor < 0{
			break;
		}else if floor == orderFloor && cart.curDir() == orderDirection{
			return 0;
		}else if cart.checkIfCommand(floor) == 1{
			turnFloor = floor;
		}
	}
	for floor := (turnFloor - cart.curDir()); floor < N_FLOORS; floor -= cart.curDir() {
		if floor < 0 {
			break;
		}else if floor == orderFloor{
			return 0;
		}else if cart.checkIfCommand(floor) == 1{
			lastFloor = floor
		}
	}
	return abs(lastFloor - orderFloor) * travelTime;
}

func costForClient(orders Orders, cart Cart, orderFloor int, orderDirection int) int {
	stops := 0;
	turnFloor := cart.curPos();
	noStops := true;

	if cart.curDir() == 0{
		return abs(orderFloor - cart.curPos()) * travelTime;
	}
	for floor := cart.curPos(); floor < N_FLOORS; floor += cart.curDir() {
		fmt.Println("floor: ", floor)
		fmt.Println("Check if command", cart.checkIfCommand(floor))
		if floor < 0{
			fmt.Println("BREAK");
			break;
		}else if floor == orderFloor && orderDirection == cart.curDir() {
			fmt.Println("number of stops: ", stops);
			return abs(cart.curPos() - floor) * travelTime + stops * stopTime;
		}else if cart.checkIfCommand(floor) == 1{
			fmt.Println("stops += 1, floor: ", floor);
			stops += 1;
			turnFloor = floor;
			noStops = false;
		}
	}
	floor := turnFloor - 1;
	for floor < N_FLOORS{
		if floor < 0 {
			fmt.Println("looped to second end floor");
			if noStops{
				return abs(cart.curPos() - orderFloor);
			}
			return (abs(cart.curPos() - turnFloor) + abs(floor - orderFloor) + N_FLOORS - 1) * travelTime + stops * stopTime;
		}else if floor == orderFloor && orderDirection == -cart.curDir(){
			return (abs(cart.curPos() - turnFloor) + abs(floor - orderFloor)) * travelTime + stops * stopTime;
		}else if cart.checkIfCommand(floor) == 1{
			stops += 1;
			noStops = false;
		}
		floor -= cart.curDir();
	}
	return 0
}

func orderIsBestForMe(orders Orders, carts []*Cart, thisCart Cart, orderFloor int, orderDirection int) int{
	lowestCost := 100000.0;
	bestCart := (*Cart)(nil)
	for cart := 0; cart < len(carts); cart++ {
		fmt.Println("Going in to orderIsBestForMe");
		fmt.Println("Checking cost of Cart: ", carts[cart]);
		clientCost := costForClient(orders, *carts[cart], orderFloor, orderDirection);
		fmt.Println("Client cost:", clientCost);
		addedCost := addedCostForElevator(orders, *carts[cart], orderFloor, orderDirection);
		fmt.Println("Added cost:", addedCost);
		cost := float64(clientCost) * 1 + float64(addedCost) * 1;
		fmt.Println("Cart number:", cart + 1, "cost: ", cost);
		if cost < lowestCost{
			lowestCost = cost;
			bestCart = carts[cart];
		}
	}
	if *bestCart == thisCart{
		return 1;
	}else{
		return 0;
	}
}

func getDirection(orders Orders, carts []*Cart, cart Cart) int {
	iterateDirection := cart.curDir();
	if iterateDirection == 0{
		iterateDirection = 1;
	}
	floor := cart.curPos();
	for floor < N_FLOORS{
		if floor < 0{
			break;
		}else if cart.checkIfCommand(floor) == 1{
			fmt.Println("Command in direction");
			return iterateDirection;
		}else if ((orders.checkIfOrder(floor, iterateDirection) == 1) && (orderIsBestForMe(orders, carts, cart, floor, iterateDirection) == 1)) || ((orders.checkIfOrder(floor, -iterateDirection) == 1) && (orderIsBestForMe(orders, carts, cart, floor, -iterateDirection) == 1)){
			fmt.Println("Found order and order is best for me, iterating in curDir");
			return iterateDirection;
		}
		floor += iterateDirection
	}
	floor = cart.curPos();
	for floor < N_FLOORS{
		if floor < 0{
			break;
		}else if cart.checkIfCommand(floor) == 1 {
			fmt.Println("Did not find orders for me or commands in dir, found command in -curDir")
			return -iterateDirection;
		}else if ((orders.checkIfOrder(floor, iterateDirection) == 1) && (orderIsBestForMe(orders, carts, cart, floor, iterateDirection) == 1)) || ((orders.checkIfOrder(floor, -iterateDirection) == 1) && (orderIsBestForMe(orders, carts, cart, floor, -iterateDirection) == 1)){
			return -iterateDirection;
		}
		floor -= iterateDirection;
	}
	return 0;
}



func main() {
	
	cart1 := Cart{dir: 1, pos: 1};
	cart2 := Cart{dir: 1, pos: 0};
	carts := []*Cart{&cart1, &cart2};
	//cart1.dir = 1;
	//cart1.pos = 1;
	//cart2.pos = 2;
	//cart2.dir = 0;
	orders.addOrder(1, -1);
	
	cart1.addCommand(0);
	cart1.addCommand(3);

	//fmt.Println("Check if command", cart1.checkIfCommand(2))
	//orders.addOrder(1,-1);
	fmt.Println("Direction for cart number 2:", getDirection(orders, carts, cart2));
}






