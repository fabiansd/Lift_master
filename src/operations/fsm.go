package operations

import (
	"driver"
	"fmt"
	"time"
)

var elevator Elevator

func Fsm_elevator() Elevator {
	return elevator
}
func Fsm_floor() int {
	return elevator.Floor
}
func Fsm_direction() Direction {
	return (elevator.Dir)
}
func Fsm_printstatus() {
	for {
		for newlines := 0; newlines < 50; newlines++ {
			fmt.Println("")
		}
		fmt.Println("Current floor: ", elevator.Floor, "\n")
		fmt.Println("Direction: ", elevator.Dir, "\n")
		fmt.Println("Requests:\n     Down | UP | Cab \n4: ", elevator.Requests[3], "\n3: ", elevator.Requests[2], "\n2: ", elevator.Requests[1], "\n1: ", elevator.Requests[0], "\n")
		switch elevator.Behaviour {
		case EB_Idle:
			fmt.Println("Elevator behaviour: IDLE\n")

		case EB_DoorOpen:
			fmt.Println("Elevator behaviour: DOOR OPEN\n")

		case EB_Moving:
			fmt.Println("Elevator behaviour: MOVING\n")

		}

		time.Sleep(300 * time.Millisecond)
	}
}

func Returnelevatorfloor() int {
	return elevator.Floor
}

func setAllLights(es Elevator) {
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			driver.Elev_set_button_lamp(btn, floor, es.Requests[floor][btn])
		}
	}
}

//Setting new order in the elevators requests
func Fsm_neworder(btn_floor int, btn_type int) {

	switch elevator.Behaviour {

	case EB_DoorOpen:
		if elevator.Floor == btn_floor {
			Timer_start()
		} else {
			elevator.Requests[btn_floor][btn_type] = true
		}
		break

	case EB_Moving:
		elevator.Requests[btn_floor][btn_type] = true
		break

	case EB_Idle:
		elevator.Requests[btn_floor][btn_type] = true

		elevator.Dir = Requests_chooseDirection(elevator)

		if elevator.Dir == DIRN_STOP {
			fmt.Println("Clear the floor!=")
			driver.Elev_set_door_open_lamp(true)
			elevator = Requests_clearAllAtCurrentFloor(elevator)
			Timer_start()
			elevator.Behaviour = EB_DoorOpen
		} else {
			driver.Elev_set_motor_direction(int(elevator.Dir))
			elevator.Behaviour = EB_Moving
		}
		break

	}

	setAllLights(elevator)
}

func Fsm_onFloorArrival(newFloor int) {

	elevator.Floor = newFloor
	driver.Elev_set_floor_indicator(elevator.Floor)

	if Requests_shouldStop(elevator) { //&& elevator.behaviour == MOVING??
		driver.Elev_set_motor_direction(int(DIRN_STOP))
		driver.Elev_set_door_open_lamp(true)
		elevator = Requests_clearAtCurrentFloor(elevator)
		Timer_start()
		setAllLights(elevator)
		elevator.Behaviour = EB_DoorOpen
	}

}

func Fsm_onDoorTimeout() {

	switch elevator.Behaviour {
	case EB_DoorOpen:
		elevator.Dir = Requests_chooseDirection(elevator)

		driver.Elev_set_door_open_lamp(false)
		driver.Elev_set_motor_direction(int(elevator.Dir))

		if elevator.Dir == DIRN_STOP {
			elevator.Behaviour = EB_Idle
		} else {
			elevator.Behaviour = EB_Moving
		}

		break
	default:
		break
	}
}
