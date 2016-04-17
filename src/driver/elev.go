package driver

import (
	"log"
)

// Wrapper for libComedi Elevator control.
// These functions provides an interface to the elevators in the real time lab

const N_FLOORS = 4
const N_BUTTONS = 3

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_DOWN1, LIGHT_UP1, LIGHT_COMMAND1},
	{LIGHT_DOWN2, LIGHT_UP2, LIGHT_COMMAND2},
	{LIGHT_DOWN3, LIGHT_UP3, LIGHT_COMMAND3},
	{LIGHT_DOWN4, LIGHT_UP4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_DOWN1, BUTTON_UP1, BUTTON_COMMAND1},
	{BUTTON_DOWN2, BUTTON_UP2, BUTTON_COMMAND2},
	{BUTTON_DOWN3, BUTTON_UP3, BUTTON_COMMAND3},
	{BUTTON_DOWN4, BUTTON_UP4, BUTTON_COMMAND4},
}

func Elev_init() bool {

	if Io_init() == false {
		return false
	}

	for floor := 0; floor < N_FLOORS; floor++ {
		for button := 0; button < N_BUTTONS; button++ {
			Elev_set_button_lamp(button, floor, false)
		}
	}
	Elev_set_door_open_lamp(false)
	Elev_set_floor_indicator(0)
	return true
}

func Elev_set_motor_direction(dirn int) {
	if dirn == 0 {
		Io_write_analog(MOTOR, 0)
	} else if dirn > 0 {
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
	} else if dirn < 0 {
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
	}
}

func Elev_set_door_open_lamp(value bool) {
	if value == true {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_get_floor_sensor_signal() int {
	if Io_read_bit(SENSOR_FLOOR1) {
		return 0
	} else if Io_read_bit(SENSOR_FLOOR2) {
		return 1
	} else if Io_read_bit(SENSOR_FLOOR3) {
		return 2
	} else if Io_read_bit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func Elev_set_floor_indicator(floor int) {
	if floor < 0 {
		log.Fatalf("Floor number is negative!")
	}
	if floor >= N_FLOORS {
		log.Fatalf("Floornumber is above topfloor")
	}

	if floor&0x02 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}

	if floor&0x01 > 0 {
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func Elev_get_button_signal(button int, floor int) bool {
	if floor >= 0 && floor < N_FLOORS {
		if button >= 0 && button < N_BUTTONS {
			return Io_read_bit(button_channel_matrix[floor][button])
		}
	}
	return false
}

func Elev_set_button_lamp(button int, floor int, value bool) {
	if floor >= 0 && floor < N_FLOORS {
		if button >= 0 && button < 3 {
			if value {
				Io_set_bit(lamp_channel_matrix[floor][button])
			} else {
				Io_clear_bit(lamp_channel_matrix[floor][button])
			}
		}
	}
}
