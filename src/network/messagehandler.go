package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"operations"
	"os"
)

func Init(outgoingMsg, incomingMsg chan operations.Udp_message) {
	// Ports randomly chosen to reduce likelihood of port collision.
	const localListenPort = 22010
	const broadcastListenPort = 22011

	const messageSize = 1024

	var udpSend = make(chan udpMessage)
	var udpReceive = make(chan udpMessage, 10)
	err := udpInit(localListenPort, broadcastListenPort, messageSize, udpSend, udpReceive)
	if err != nil {
		fmt.Print("UdpInit() error: %v \n", err)
	}

	go forwardOutgoing(outgoingMsg, udpSend)
	go forwardIncoming(incomingMsg, udpReceive)

	log.Println("Network initialised.")
}

// aliveSpammer periodically sends messages on the network to notify it is alive
func GetLocalIP() string {
	var localIP string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}
	return localIP
}

// forwardOutgoing continuosly checks for messages to be sent on the network
// by reading the OutgoingMsg channel. Each message read is sent to the udp file
// as json.
func forwardOutgoing(outgoingMsg <-chan operations.Udp_message, udpSend chan<- udpMessage) {
	for {
		msg := <-outgoingMsg

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("%sjson.Marshal error: %v\n%s")
		}

		udpSend <- udpMessage{raddr: "broadcast", data: jsonMsg, length: len(jsonMsg)}
	}
}

func forwardIncoming(incomingMsg chan<- operations.Udp_message, udpReceive <-chan udpMessage) {
	for {
		udpMessage := <-udpReceive
		var message operations.Udp_message

		if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
			fmt.Printf("json.Unmarshal error: %s\n", err)
		}

		message.Addr = udpMessage.raddr
		incomingMsg <- message
	}
}
