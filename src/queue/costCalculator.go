package queue

import (
	"driver"
	"log"
	"operations"
)

//Calculates cost by simulating the new order among the existing local orders for the elevator
//Each floor-stop and travel between floors has cost of 2
//If the elevator starts between two floors, the cost will only be 1 for this
func CalCost(tarFloor, tarBtn, prevFloor, curFloor, curDir int) int {
	//fmt.Println("copy of elevator")
	var qCopy operations.Elevator
	qCopy.Floor = operations.Fsm_floor()     //Copies floor
	qCopy.Dir = operations.Direction(curDir) //Copies the direction of the local elevator

	for f := 0; f < driver.N_FLOORS; f++ { //copying local requests
		for b := 0; b < driver.N_BUTTONS; b++ {
			qCopy.Requests[f][b] = operations.Fsm_elevator().Requests[f][b]
		}
	}

	cost := 0

	qCopy.Requests[qCopy.Floor][tarBtn] = true
	//start the elevator in the right direction
	if tarFloor-curFloor < 0 {
		qCopy.Dir = -1
	} else if tarFloor-curFloor > 0 {
		qCopy.Dir = 1
	}

	if curFloor == -1 {
		//cost is set to one if elev between two floors
		cost++
	} else if qCopy.Dir != operations.DIRN_STOP {
		//cost is set to two if elev moving at floor
		cost += 2
	}

	qCopy.Floor, qCopy.Dir = incrementFloor(qCopy.Floor, int(qCopy.Dir)) //First incrementation
	//Iterates through the orders, with an upper limit of 10
	for n := 0; !(qCopy.Floor == tarFloor && operations.Requests_shouldStop(qCopy)) && n < 10; n++ {
		if operations.Requests_shouldStop(qCopy) {
			cost += 2
			qCopy = operations.Requests_clearAtCurrentFloor(qCopy,true)

		}
		qCopy.Dir = operations.Requests_chooseDirection(qCopy)
		qCopy.Floor, qCopy.Dir = incrementFloor(qCopy.Floor, int(qCopy.Dir))
		cost += 2
		if !(operations.Requests_below(qCopy) && operations.Requests_above(qCopy)) {
			break
		}
	}
	return cost
}

func incrementFloor(floor, dir int) (int, operations.Direction) {
	switch dir {
	case operations.DIRN_DOWN:
		floor--
	case operations.DIRN_UP:
		floor++
	case operations.DIRN_STOP:
		//nun
	default:
		log.Fatalln("incrementFloor(): invalid dir, not incremented")
	}
	if floor <= 0 && dir == operations.DIRN_DOWN {
		dir = operations.DIRN_UP
		floor = 0
	}
	if floor >= driver.N_FLOORS-1 && dir == operations.DIRN_UP {
		dir = operations.DIRN_DOWN
		floor = driver.N_FLOORS - 1
	}
	return floor, operations.Direction(dir)
}
