package globalOperations

import (
	"driver"
	"elevatorOperations"
	"log"
)

//Calculates cost by simulating the new order among the existing local orders for the elevator
//Each floor-stop and travel between floors has cost of 2
//If the elevator starts between two floors, the cost will only be 1 for this
func CalCost(tarFloor, tarBtn, prevFloor, curFloor, curDir int) int {
	//fmt.Println("copy of elevator")
	var qCopy elevatorOperations.Elevator
	qCopy.Floor = elevatorOperations.Fsm_floor()     //Copies floor
	qCopy.Dir = elevatorOperations.Direction(curDir) //Copies the direction of the local elevator

	for f := 0; f < driver.N_FLOORS; f++ { //copying local requests
		for b := 0; b < driver.N_BUTTONS; b++ {
			qCopy.Requests[f][b] = elevatorOperations.Fsm_elevator().Requests[f][b]
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
	} else if qCopy.Dir != elevatorOperations.DIRN_STOP {
		//cost is set to two if elev moving at floor
		cost += 2
	}

	qCopy.Floor, qCopy.Dir = incrementFloor(qCopy.Floor, int(qCopy.Dir)) //First incrementation
	//Iterates through the orders, with an upper limit of 10
	for n := 0; !(qCopy.Floor == tarFloor && elevatorOperations.Requests_shouldStop(qCopy)) && n < 10; n++ {
		if elevatorOperations.Requests_shouldStop(qCopy) {
			cost += 2
			qCopy = elevatorOperations.Requests_clearAtCurrentFloor(qCopy, true)

		}
		qCopy.Dir = elevatorOperations.Requests_chooseDirection(qCopy)
		qCopy.Floor, qCopy.Dir = incrementFloor(qCopy.Floor, int(qCopy.Dir))
		cost += 2
		if !(elevatorOperations.Requests_below(qCopy) && elevatorOperations.Requests_above(qCopy)) {
			break
		}
	}
	return cost
}

func incrementFloor(floor, dir int) (int, elevatorOperations.Direction) {
	switch dir {
	case elevatorOperations.DIRN_DOWN:
		floor--
	case elevatorOperations.DIRN_UP:
		floor++
	case elevatorOperations.DIRN_STOP:
		//nun
	default:
		log.Fatalln("incrementFloor(): invalid dir, not incremented")
	}
	if floor <= 0 && dir == elevatorOperations.DIRN_DOWN {
		dir = elevatorOperations.DIRN_UP
		floor = 0
	}
	if floor >= driver.N_FLOORS-1 && dir == elevatorOperations.DIRN_UP {
		dir = elevatorOperations.DIRN_DOWN
		floor = driver.N_FLOORS - 1
	}
	return floor, elevatorOperations.Direction(dir)
}
