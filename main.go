package main

import (
	"driver"
	"network"
	"operations"
	"queue"
	"fmt"
	"os/signal"
	"os"
	//"time"
	"log"
)
//control channels
var outgoingMsg = make(chan operations.Udp_message, 10)
var incomingMsg = make(chan operations.Udp_message, 10)
var newOrderChan = make(chan operations.Keypress)
var costChan = make(chan operations.Udp_message)


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
	go terminateEngine()
	go networkHandler()
	
	network.Init(outgoingMsg, incomingMsg)	

	fmt.Println("Start listening")

	outgoingMsg <- operations.Udp_message{Category: operations.NewOrder, Floor: -1, Button: -1, Cost: 10000}
	
	//operations.CloseConnectionChan <- true


	//cmd := exec.Command("gnome-terminal", "-x", "go", "run", "main.go")
	//cmd.Run()
	
	//msg := network.Udp_message{Raddr: "broadcast", Data: object.Data, Length: 0}
	halt := make(chan bool)
	<- halt

}

func Pollneworder(){
		for{
			msg := <-newOrderChan

			switch msg.Button{
			case operations.B_Inside:
				fmt.Println(msg.Button)
				operations.Fsm_neworder(msg.Floor, msg.Button)
			case operations.B_Up:
				outgoingMsg <- operations.Udp_message{Category: operations.NewOrder, Floor: msg.Floor, Button: msg.Button, Cost: -1}
			case operations.B_Down:
				outgoingMsg <- operations.Udp_message{Category: operations.NewOrder, Floor: msg.Floor, Button: msg.Button, Cost: -1}
			default:
				fmt.Println("Error: non-buttontype")
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


func networkHandler(){
	for{
		select{
		case msg := <- incomingMsg:
			messageHandler(msg)
		}
	}
}

func messageHandler(message operations.Udp_message){
	switch message.Category{

	case operations.Cost:
		costChan <- message
	case operations.Livefeed:
		fmt.Println("recieving livefeed")
	case operations.NewOrder:
		
		var b int = message.Button //KLARER IKKE BRUKE OVERFØRT INT I NEWORDERS FUNKSJONEN
		var a int
		switch b{
		case 0:
			a = 0
		case 1:
			a = 1
		}
		fmt.Println("neworder - call the cost!")
		cost := queue.CalCost(message.Floor, message.Button, operations.Fsm_elevator().PrevFloor, operations.Fsm_elevator().Floor, a)
		outgoingMsg <- operations.Udp_message{Category: operations.Cost, Floor: message.Floor, Button: message.Button, Cost: cost}
	case operations.CompletedOrder:
	default:
	}
}

func terminateEngine() {//kills engine when program is termianted
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.Elev_set_motor_direction(operations.DIRN_STOP)
	log.Fatal("User terminated program")
}