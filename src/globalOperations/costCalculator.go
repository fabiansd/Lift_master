package globalOperations

import (
	"driver"
	"elevatorOperations"
	"log"
)

//Calculates cost by simulating the new order appended to the local elevator 
func CalCost(tarFloor, tarBtn, prevFloor, curFloor, curDir int) int {
	var qCopy elevatorOperations.Elevator
	qCopy.Floor = elevatorOperations.Fsm_floor()
	qCopy.Dir = elevatorOperations.Direction(curDir)

	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			qCopy.Requests[f][b] = elevatorOperations.Fsm_elevator().Requests[f][b]
		}
	}
	cost := 0
	qCopy.Requests[qCopy.Floor][tarBtn] = true

	if tarFloor-curFloor < 0 {
		qCopy.Dir = -1
	} else if tarFloor-curFloor > 0 {
		qCopy.Dir = 1
	}
	// local elevator copied

	if curFloor == -1 {
		cost++
	} else if qCopy.Dir != elevatorOperations.DIRN_STOP {
		cost += 2
	}
	
	// Simulation iteration-cap set to 10
	//Each floor-stop and travel between floors has cost of 2
	//If the elevator starts between two floors, the cost will be 1
	qCopy.Floor, qCopy.Dir = incrementFloor(qCopy.Floor, int(qCopy.Dir))
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
		//none
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
