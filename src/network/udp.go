package network

import (
	"elevatorOperations"
	"net"
	"strconv"
	"time"
)

type UdpConnection struct {
	Addr  string
	Timer *time.Timer
}

var baddr *net.UDPAddr

type udpMessage struct {
	raddr  string
	data   []byte
	length int
}

func udpInit(localListenPort, broadcastListenPort, message_size int, send_ch, receive_ch chan udpMessage) (err error) {
	baddr, err = net.ResolveUDPAddr("udp4", "129.241.187.255:"+strconv.Itoa(broadcastListenPort))
	if err != nil {
		return err
	}

	tempConn, err := net.DialUDP("udp4", nil, baddr)
	defer tempConn.Close()
	tempAddr := tempConn.LocalAddr()
	laddr, err := net.ResolveUDPAddr("udp4", tempAddr.String())
	laddr.Port = localListenPort
	elevatorOperations.Laddr = laddr.String()

	localListenConn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return err
	}

	broadcastListenConn, err := net.ListenUDP("udp", baddr)
	if err != nil {
		localListenConn.Close()
		return err
	}

	go udp_receive_server(localListenConn, broadcastListenConn, message_size, receive_ch)
	go udp_transmit_server(localListenConn, broadcastListenConn, send_ch)

	return err
}

func udp_transmit_server(lconn, bconn *net.UDPConn, send_ch <-chan udpMessage) {

	for {
		msg := <-send_ch
		if msg.raddr == "broadcast" {
			_, _ = lconn.WriteToUDP(msg.data, baddr)
		} else {
			raddr, _ := net.ResolveUDPAddr("udp", msg.raddr)
			_, _ = lconn.WriteToUDP(msg.data, raddr)
		}
	}
}

func udp_receive_server(lconn, bconn *net.UDPConn, message_size int, receive_ch chan<- udpMessage) {

	bconn_rcv_ch := make(chan udpMessage)
	lconn_rcv_ch := make(chan udpMessage)

	go udp_connection_reader(lconn, message_size, lconn_rcv_ch)
	go udp_connection_reader(bconn, message_size, bconn_rcv_ch)

	for {
		select {

		case buf := <-bconn_rcv_ch:
			receive_ch <- buf

		case buf := <-lconn_rcv_ch:
			receive_ch <- buf
		}
	}
}

func udp_connection_reader(conn *net.UDPConn, message_size int, rcv_ch chan<- udpMessage) {

	for {
		buf := make([]byte, message_size)
		n, raddr, _ := conn.ReadFromUDP(buf)
		buf = buf[:n]
		rcv_ch <- udpMessage{raddr: raddr.String(), data: buf, length: n}
	}
}
