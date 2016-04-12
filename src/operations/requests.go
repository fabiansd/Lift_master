package operations

import (
	"driver"
	"fmt"
	"time"
)

func Requests_above(e Elevator) bool {
	for f := e.Floor + 1; f < driver.N_FLOORS; f++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if Fsm_elevator().Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_below(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if Fsm_elevator().Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_chooseDirection(e Elevator) Direction {
	switch e.Dir {
	case DIRN_UP:
		if Requests_above(e) {
			return DIRN_UP
		} else if Requests_below(e) {
			return DIRN_DOWN
		} else {
			return DIRN_STOP
		}

	case DIRN_DOWN:
		if Requests_below(e) {
			return DIRN_DOWN
		} else if Requests_above(e) {
			return DIRN_UP
		} else {
			return DIRN_STOP
		}
	case DIRN_STOP:
		if Requests_above(e) {
			return DIRN_UP
		} else if Requests_below(e) {
			return DIRN_DOWN
		} else {
			return DIRN_STOP
		}
	default:
		return DIRN_STOP
	}
}

func Requests_shouldStop(e Elevator) bool {
	fmt.Println("Request_ss: q1 ")
	fmt.Println(e)
	switch e.Dir {
	case DIRN_DOWN:
		return (e.Requests[e.Floor][0] || e.Requests[e.Floor][2] || (!Requests_below(e)))
	case DIRN_UP:
		return (e.Requests[e.Floor][1] || e.Requests[e.Floor][2] || (!Requests_above(e)))
	case DIRN_STOP:
		return true
	default:
		return true
	}
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {

	e.Requests[e.Floor][B_Inside] = false
	switch e.Dir {
	case DIRN_UP:
		e.Requests[e.Floor][B_Up] = false
		if !Requests_above(e) {
			e.Requests[e.Floor][B_Down] = false
		}
		break

	case DIRN_DOWN:
		e.Requests[e.Floor][B_Down] = false
		if !Requests_below(e) {
			e.Requests[e.Floor][B_Up] = false
		}
		break

	case DIRN_STOP:
	default:
		e.Requests[e.Floor][B_Up] = false
		e.Requests[e.Floor][B_Down] = false
		break
	}
	return e
}

func Request_buttons(newOrderChan chan Keypress) {

	for {
		var prevReq = [driver.N_FLOORS][driver.N_BUTTONS]bool{}

		for floor := 0; floor < driver.N_FLOORS; floor++ {
			for btn := 0; btn < driver.N_BUTTONS; btn++ {
				buttonPressed := driver.Elev_get_button_signal(btn, floor)
				if buttonPressed && buttonPressed != prevReq[floor][btn] {
					//Fsm_onRequestButtonPress(floor, ButtonType(btn),newOrderChan)
					newOrderChan <- Keypress{floor, btn}
				}
				prevReq[floor][btn] = buttonPressed
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func Request_floorSensor() {
	for {
		var prevFloor = elevator.Floor
		floorSensor := driver.Elev_get_floor_sensor_signal()
		if floorSensor != -1 && floorSensor != prevFloor {
			Fsm_onFloorArrival(floorSensor)
		}
		prevFloor = floorSensor
		elevator.PrevFloor = prevFloor
		time.Sleep(25 * time.Millisecond)
	}
}
func Request_timecheck() {
	for {
		if Timer_timedout() {
			Fsm_onDoorTimeout()
			Timer_stop()
		}
		time.Sleep(25 * time.Millisecond)
	}
}
