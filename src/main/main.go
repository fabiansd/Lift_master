package main

import (
	"driver"
	"elevatorOperations"
	"fmt"
	"globalOperations"
	"log"
	"network"
	"os"
	"os/signal"
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

	//HAII
	//Control channels
	var outgoingMsg = make(chan elevatorOperations.Udp_message, 1000)
	var incomingMsg = make(chan elevatorOperations.Udp_message, 1000)
	var newOrderChan = make(chan elevatorOperations.Keypress, 100)
	var costChan = make(chan elevatorOperations.Udp_message, 100)
	var aliveChannel = make(chan []string)
	var elevatorsOnlineChan = make(chan int)
	var orderCompleteChannel = make(chan globalOperations.Order, 100)
	var orderDeleted = make(chan elevatorOperations.Keypress, 100)

	Laddr := network.GetLocalIP() + ":22010"
	fmt.Println("Laddr initialized", Laddr)
	network.Init(outgoingMsg, incomingMsg)

	go network.UdpSendAlive("30000")
	go network.UdpRecvAlive("30000", aliveChannel)
	go elevatorOperations.Request_Poll(orderDeleted, newOrderChan)
	go terminateEngine()
	go globalOperations.BackupHandler(newOrderChan, orderCompleteChannel, outgoingMsg)
	go globalOperations.RecieveCosts(costChan, outgoingMsg, elevatorsOnlineChan)

	//Handle events
	for {
		select {
		//Poll for new messages from the UDP network
		case message := <-incomingMsg:
			switch message.Category {
			case elevatorOperations.Cost:
				costChan <- message
			case elevatorOperations.Killfeed:
				if Laddr == message.AssignAddr {

				}
			case elevatorOperations.NewOrder:
				elevatorOperations.SetGlobalLights(message.Floor, message.Button, true)
				cost := globalOperations.CalCost((message.Floor), (message.Button), elevatorOperations.Fsm_floor(), driver.Elev_get_floor_sensor_signal(), int(elevatorOperations.Fsm_direction()))
				outgoingMsg <- elevatorOperations.Udp_message{Category: elevatorOperations.Cost, Floor: message.Floor, Button: message.Button, Cost: cost}
			case elevatorOperations.CompletedOrder:
				elevatorOperations.SetGlobalLights(message.Floor, message.Button, false)
				orderCompleteChannel <- globalOperations.Order{Floor: message.Floor, Button: message.Button}

			case elevatorOperations.AssignedOrder:
				if Laddr == message.AssignAddr {
					elevatorOperations.Fsm_neworder(message.Floor, message.Button)
				}
			default:
			}
		//Poll for new orders
		case neworder := <-newOrderChan:
			switch neworder.Button {
			case elevatorOperations.B_Inside:
				elevatorOperations.Fsm_neworder(neworder.Floor, neworder.Button)
			case elevatorOperations.B_Up:
				outgoingMsg <- elevatorOperations.Udp_message{Category: elevatorOperations.NewOrder, Floor: (neworder.Floor), Button: neworder.Button, Cost: 0}
			case elevatorOperations.B_Down:
				outgoingMsg <- elevatorOperations.Udp_message{Category: elevatorOperations.NewOrder, Floor: (neworder.Floor), Button: neworder.Button, Cost: 0}
			default:
				fmt.Println("Error: non-buttontype")
			}
		case alive := <-aliveChannel:
			elevatorsOnline := len(alive)
			fmt.Println(elevatorOperations.Yellow, "Number of alive elevators set to: ", elevatorsOnline, elevatorOperations.White)
			elevatorsOnlineChan <- elevatorsOnline
		case orderDeleted := <-orderDeleted:

			outgoingMsg <- elevatorOperations.Udp_message{Category: elevatorOperations.CompletedOrder, Floor: orderDeleted.Floor, Button: orderDeleted.Button, Cost: 0}
		}
	}
}

func terminateEngine() { //kills engine when program is termianted
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.Elev_set_motor_direction(elevatorOperations.DIRN_STOP)
	log.Fatal("User terminated program")
}
