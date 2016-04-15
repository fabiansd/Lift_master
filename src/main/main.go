package main

import (
	"driver"
	"fmt"
	"network"
	"operations"
	"os"
	"os/signal"
	"queue"
	"log"
)
//Global channels

var newOrderChan = make(chan operations.Keypress, 100)


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
	//Control channels
	var outgoingMsg = make(chan operations.Udp_message, 1000)
	var incomingMsg = make(chan operations.Udp_message, 1000)
	var costChan = make(chan operations.Udp_message, 100)
	var aliveChannel = make(chan []string)
	var elevatorsOnlineChan = make(chan int)
	var orderCompleteChannel = make(chan queue.Order, 100)
	var orderDeleted = make(chan operations.Keypress, 100)

	Laddr := network.GetLocalIP() + ":22010"
	fmt.Println("Laddr initialized", Laddr)
	network.Init(outgoingMsg, incomingMsg)

	go network.UdpSendAlive("30000")
	go network.UdpRecvAlive("30000", aliveChannel)
	go operations.Request_Poll(orderDeleted)
	go operations.Request_buttons(newOrderChan)
	go operations.Request_floorSensor()
	go terminateEngine()
	//go operations.Fsm_printstatus()
	go operations.Request_timecheck()
	go queue.BackupHandler(newOrderChan,orderCompleteChannel,outgoingMsg)
	go queue.RecieveCosts(costChan, outgoingMsg, elevatorsOnlineChan)

	//Handle events
	for {
		select {
		//Poll for new messages from the UDP network
		case message := <-incomingMsg:
			switch message.Category {
			case operations.Cost:
				costChan <- message
			case operations.Killfeed:
				if Laddr == message.AssignAddr{
					fmt.Println(operations.Red, "ERROR: System function failure", operations.Red)
					fmt.Println(operations.White)
					driver.Elev_set_motor_direction(operations.DIRN_STOP)
					os.Exit(1)
				}
				//fmt.Println("recieving livefeed")
			case operations.NewOrder:
				operations.SetGlobalLights(message.Floor,message.Button,true)
				cost := queue.CalCost((message.Floor), (message.Button), operations.Fsm_floor(), driver.Elev_get_floor_sensor_signal(), int(operations.Fsm_direction()))
				outgoingMsg <- operations.Udp_message{Category: operations.Cost, Floor: message.Floor, Button: message.Button, Cost: cost}
			case operations.CompletedOrder:
				//Send the completed order to orderCompleteChannel to registrer the completion in backup
				operations.SetGlobalLights(message.Floor,message.Button,false)
				orderCompleteChannel <- queue.Order{Floor: message.Floor, Button: message.Button}

			case operations.AssignedOrder:
				if Laddr == message.AssignAddr {
					operations.Fsm_neworder(message.Floor, message.Button)
				}
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
		case alive := <-aliveChannel:
			elevatorsOnline := len(alive)
			fmt.Println("Number of alive elevators has been changed to: ", elevatorsOnline)
			elevatorsOnlineChan <- elevatorsOnline
		case orderDeleted:= <-orderDeleted:
			
			outgoingMsg <- operations.Udp_message{Category: operations.CompletedOrder, Floor: orderDeleted.Floor, Button: orderDeleted.Button, Cost: 0}
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
