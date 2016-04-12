package queue

import (
	"driver"
	"fmt"
	"log"
	"operations"
)

func CalCost(tarFloor, tarBtn, prevFloor, curFloor, curDir int) int {
	fmt.Println("copy of elevator")
	var q1 operations.Elevator
	q1.Floor = 1
	q1.PrevFloor = 1
	q1.Dir = operations.Direction(1)
	q1.Behaviour = operations.EB_Moving
	//q1.Requests = operations.[driver.N_FLOORS][driver.N_BUTTONS]bool

	//var q operations.Elevator
	//fmt.Println("new")
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			q1.Requests[f][b] = operations.Fsm_elevator().Requests[f][b]
		}
	}

	cost := 0
	//q.Floor = prevFloor
	//q.Dir = operations.Direction(curDir)
	//q.Requests[q.Floor][operations.B_Inside] = true
	//fmt.Println("copied")
	fmt.Println(q1)

	if curFloor == -1 {
		//cost is set to one if elev between two floors
		cost++
	} else if q1.Dir != operations.DIRN_STOP {
		cost += 2
		//cost is set to two if elev moving at floor
	}
	q1.Floor, q1.Dir = incrementFloor(q1.Floor, int(q1.Dir))
	//fmt.Println(cost)
	for n := 0; !(q1.Floor == tarFloor && operations.Requests_shouldStop(q1)) && n < 10; n++ {
		//if !(q1.Floor == tarFloor) &&

		//!(q1.Floor == tarFloor && operations.Requests_shouldStop(q1)) && {
		if operations.Requests_shouldStop(q1) {
			cost += 2
			q1.Requests[q1.Floor][operations.B_Up] = false
			q1.Requests[q1.Floor][operations.B_Down] = false
			q1.Requests[q1.Floor][operations.B_Inside] = false
		}
		q1.Dir = operations.Requests_chooseDirection(q1)
		q1.Floor, q1.Dir = incrementFloor(q1.Floor, int(q1.Dir))
		cost += 2
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
		//No in/decrement
	default:
		//CloseConnectionChan <- true
		//Restart.Run()
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
