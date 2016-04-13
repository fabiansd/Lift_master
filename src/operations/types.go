package operations

import (
	"driver"
)

//Types og directions
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

var Laddr string //local IP address
var CloseConnectionChan = make(chan bool)

type Keypress struct {
	Floor  int
	Button int
}

const (
	Livefeed int = iota + 1
	NewOrder
	CompletedOrder
	Cost
	AssignedOrder
)
