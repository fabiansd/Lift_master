package main

import (
	"driver"
	"fmt"
	"network"
	"operations"
	"os"
	"os/signal"
	"queue"
	//"time"
	"log"
)

func main() {
	//Initializes the elevator
	driver.Elev_init()
	driver.Elev_set_motor_direction(-1)
	for {
		if driver.Elev_get_floor_sensor_signal() != -1 {
			driver.Elev_set_motor_direction(0)
			driver.Elev_set_floor_indicator(driver.Elev_get_floor_sensor_signal())

			break
		}
	}
	fmt.Println(operations.Fsm_elevator())
	//operations.Laddr = "129.241.187.158"
	//Control channels
	var outgoingMsg = make(chan operations.Udp_message, 1000)
	var incomingMsg = make(chan operations.Udp_message, 1000)
	var newOrderChan = make(chan operations.Keypress, 100)
	var costChan = make(chan operations.Udp_message, 100)

	network.Init(outgoingMsg, incomingMsg)
	go operations.Request_buttons(newOrderChan)
	go operations.Request_floorSensor()
	go terminateEngine()
	//go operations.Fsm_printstatus()
	go operations.Request_timecheck()
	go queue.RecieveCosts(costChan)
	fmt.Println("Startup done")

	//Handle events
	for {
		select {
		//Poll for new messages from the UDP network
		case message := <-incomingMsg:
			switch message.Category {
			case operations.Cost:
				costChan <- message
			case operations.Livefeed:
				//fmt.Println("recieving livefeed")
			case operations.NewOrder:
				//fmt.Println("neworder - call the cost!")
				//fmt.Println(message)
				cost := queue.CalCost((message.Floor), (message.Button), operations.Fsm_floor(), driver.Elev_get_floor_sensor_signal(), operations.Fsm_direction())
				//fmt.Println(cost)
				outgoingMsg <- operations.Udp_message{Category: operations.Cost, Floor: message.Floor, Button: message.Button, Cost: cost}
			case operations.CompletedOrder:
			default:
			}
		//Poll for new orders
		case neworder := <-newOrderChan:
			switch neworder.Button {
			case operations.B_Inside:
				operations.Fsm_neworder(neworder.Floor, neworder.Button)
			case operations.B_Up:
				outgoingMsg <- operations.Udp_message{Category: operations.NewOrder, Floor: (neworder.Floor), Button: neworder.Button, Cost: 0}
			case operations.B_Down:
				outgoingMsg <- operations.Udp_message{Category: operations.NewOrder, Floor: (neworder.Floor), Button: neworder.Button, Cost: 0}
			default:
				fmt.Println("Error: non-buttontype")
			}
		}
	}

	//outgoingMsg <- operations.Udp_message{Category: operations.NewOrder, Floor: 1, Button: 1, Cost: 10000}

	//operations.CloseConnectionChan <- true

	//cmd := exec.Command("gnome-terminal", "-x", "go", "run", "main.go")
	//cmd.Run()

	//msg := network.Udp_message{Raddr: "broadcast", Data: object.Data, Length: 0}
	halt := make(chan bool)
	<-halt

}

func terminateEngine() { //kills engine when program is termianted
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.Elev_set_motor_direction(operations.DIRN_STOP)
	log.Fatal("User terminated program")
}
