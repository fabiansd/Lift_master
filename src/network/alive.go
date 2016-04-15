package network

import (
	"net"
	"time"
)

const aliveSendInterval = 70 * time.Millisecond
const aliveTimeout = 300 * time.Millisecond

func UdpSendAlive(port string) {

	udpAddr, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:"+port)
	udpConn, _ := net.DialUDP("udp4", nil, udpAddr)

	for {
		time.Sleep(aliveSendInterval)
		udpConn.Write([]byte("I am alive"))
	}
}

func UdpRecvAlive(port string, peerListLocalCh chan []string) {

	var buf [1024]byte

	lastSeen := make(map[string]time.Time)
	hasChanges := false
	var peerList []string

	service := ":" + port
	udpAddr, _ := net.ResolveUDPAddr("udp4", service)
	readConn, _ := net.ListenUDP("udp4", udpAddr)

	for {
		hasChanges = false

		readConn.SetReadDeadline(time.Now().Add(aliveTimeout))
		_, fromAddress, err := readConn.ReadFromUDP(buf[0:])

		if err != nil {
			continue
		}

		addrString := fromAddress.IP.String()

		_, addrIsInList := lastSeen[addrString]

		if !addrIsInList {
			hasChanges = true
		}

		lastSeen[addrString] = time.Now()

		for k, v := range lastSeen {
			if time.Now().Sub(v) > aliveTimeout {
				hasChanges = true
				delete(lastSeen, k)
			}
		}

		if hasChanges {
			peerList = nil

			for k, _ := range lastSeen {
				peerList = append(peerList, k)
			}
			peerListLocalCh <- peerList
		}
	}
}
