package elevatorOperations

import (
	"driver"
)

type Direction int

const (
	DIRN_DOWN = -1 + iota
	DIRN_STOP
	DIRN_UP
)

type ButtonType int

const (
	B_Down = 0 + iota
	B_Up
	B_Inside
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = 0 + iota
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	Floor     int
	PrevFloor int
	Dir       Direction
	Behaviour ElevatorBehaviour
	Requests  [driver.N_FLOORS][driver.N_BUTTONS]bool
}

type Udp_message struct {
	Category   int
	Floor      int
	Button     int
	Cost       int
	Addr       string
	AssignAddr string
}

var Laddr string

type Keypress struct {
	Floor  int
	Button int
}

const (
	NewOrder int = iota + 1
	CompletedOrder
	Cost
	AssignedOrder
)
const White = "\x1b[37;1m"
const Red = "\x1b[31;1m"
const Yellow = "\x1b[33;1m"
