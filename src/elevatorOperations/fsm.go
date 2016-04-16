package elevatorOperations

import (
	"driver"
	"fmt"
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
func Fsm_behaviour() ElevatorBehaviour {
	return (elevator.Behaviour)
}
func Fsm_requests() [driver.N_FLOORS][driver.N_BUTTONS]bool {
	return elevator.Requests
}
func Fsm_setRequest(i int, B_Inside int) {
	elevator.Requests[i][B_Inside] = true
}

func Fsm_printrequest(temp [driver.N_FLOORS][driver.N_BUTTONS]bool, requestType string) {
	fmt.Println("\n", requestType, "\n")
	fmt.Println("     Down | UP | Cab \n4: ", temp[3], "\n3: ", temp[2], "\n2: ", temp[1], "\n1: ", temp[0], "\n")

}

func setAllLights(es Elevator) {
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			driver.Elev_set_button_lamp(btn, floor, es.Requests[floor][btn])
		}
	}
}

func SetGlobalLights(floor int, button int, value bool) {
	driver.Elev_set_button_lamp(button, floor, value)
}

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
	TakeBackup()
	setAllLights(elevator)
}

func Fsm_onFloorArrival(newFloor int) {
	elevator.Floor = newFloor
	driver.Elev_set_floor_indicator(elevator.Floor)
	if Requests_shouldStop(elevator) {
		driver.Elev_set_motor_direction(int(DIRN_STOP))
		driver.Elev_set_door_open_lamp(true)
		elevator = Requests_clearAtCurrentFloor(elevator, false)
		Timer_start()
		setAllLights(elevator)
		elevator.Behaviour = EB_DoorOpen
	}

}

func Fsm_onDoorTimeout() {

	Fsm_printrequest(Fsm_requests(), "Local elevator requests")
	if !Requests_below(elevator) {
		orderDeletedLocally <- Keypress{Floor: elevator.Floor, Button: int(B_Up)}
	}
	if !Requests_above(elevator) {
		orderDeletedLocally <- Keypress{Floor: elevator.Floor, Button: int(B_Down)}
	}
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
