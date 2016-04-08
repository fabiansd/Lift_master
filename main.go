package main

import (
	"driver"
	"network"
	"operations"
	"fmt"
	//"os/exec"
	//"time"
)
//control channels
var outgoingMsg = make(chan operations.Udp_message, 10)
var incomingMsg = make(chan operations.Udp_message, 10)
var newOrderChan = make(chan operations.Keypress)


func Pollneworder(){
	for{
		select{
		case msg := <-newOrderChan:
			fmt.Println("new order to orderchan")
			operations.Fsm_neworder(msg.Floor, msg.Button)
		default:
	}
}
}


func Initialize() bool {
	driver.Elev_init()
	driver.Elev_set_motor_direction(-1)
	for {
		if driver.Elev_get_floor_sensor_signal() != -1 {
			driver.Elev_set_motor_direction(0)
			driver.Elev_set_floor_indicator(driver.Elev_get_floor_sensor_signal())
			return true
		}
	}
}

func main() {
	if Initialize() {
		fmt.Printf("Started!\n")
	} else {
		fmt.Printf("error!\n")
	}

	operations.Laddr = "129.241.187.158"

	go operations.Request_buttons(newOrderChan)
	go operations.Request_floorSensor()
	//go operations.Fsm_printstatus()
	go operations.Request_timecheck()
	go Pollneworder()


	object := operations.Udp_message{Category: 0, Floor: 0, Button: 0, Cost: 0}
	
	network.Init(outgoingMsg, incomingMsg)	

	fmt.Println("Start listening")
	
	for {
		select{
		case msg := <- incomingMsg:
			object = msg
			fmt.Println("recieving")
			fmt.Println(object.Floor)
		}
	}
	
	//operations.CloseConnectionChan <- true


	//cmd := exec.Command("gnome-terminal", "-x", "go", "run", "main.go")
	//cmd.Run()
	
	//msg := network.Udp_message{Raddr: "broadcast", Data: object.Data, Length: 0}

}