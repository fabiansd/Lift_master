package queue

import (
	 "operations"
	 "driver"
	 "log"
 )


 func CalCost(tarFloor, tarBtn, prevFloor, curFloor, curDir int) int{
 	q := deepCopy()
 	setOrder(&q, tarFloor, operations.B_Inside,true)

 	cost := 0
 	floor := prevFloor
 	dir := curDir

 	if curFloor == -1{
 		//cost is set to one if elev between two floors
 		cost++
 	} else if dir != operations.DIRN_STOP{
 		cost += 2
 		//cost is set to two if elev moving at floor
 	}
 	floor, dir = incrementFloor(floor, dir)

 	for n :=0; !(floor == tarFloor && operations.Requests_shouldStop(q)) && n < 10; n++{
 		if operations.Requests_shouldStop(q){
 			cost +=2
 			setOrder(&q, floor, operations.B_Up, false)
 			setOrder(&q, floor, operations.B_Down, false)
 			setOrder(&q, floor, operations.B_Inside, false)
 		}
 		dir = int(operations.Requests_chooseDirection(q))
 		floor, dir = incrementFloor(floor, dir)
 		cost += 2
 	}
 	return cost
 }

 func setOrder(q *operations.Elevator,floor, button int, status bool){
 	if q.Requests[floor][button] == true{
 		return
 	}else{
 		q.Requests[floor][button] = status
	}
}


 func incrementFloor(floor, dir int) (int, int) {
 	switch dir{
 	case operations.DIRN_DOWN:
 		floor --
 	case operations.DIRN_UP:
 		floor++
 	case operations.DIRN_STOP:
 		//No in/decrement
 	default:
 		//CloseConnectionChan <- true
 		//Restart.Run()
 		log.Fatalln("incrementFloor(): invalid dir, not incremented")
 	}
 	if floor <= 0 && dir == operations.DIRN_DOWN{
 		dir = operations.DIRN_DOWN
 		floor = 0
 	}
 	if floor >= driver.N_FLOORS -1 && dir == operations.DIRN_UP{
 		dir = operations.DIRN_DOWN
 		floor = driver.N_FLOORS -1
 	}
 	return floor, dir
 	
 }


 func deepCopy() operations.Elevator {
 	elevCopy := new(operations.Elevator)
	for f := 0; f < driver.N_FLOORS; f++ {
		for b := 0; b < driver.N_BUTTONS; b++ {
			elevCopy.Requests[f][b] = operations.Fsm_elevator().Requests[f][b]
		}
	}
	return *elevCopy
}
