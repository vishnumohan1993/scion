// Code references : https://github.com/netsec-ethz/scion-homeworks/blob/master/bottleneck_bw_est/v1_bw_est_client.go 

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sciond"
)

func check(e1 error) {
	if e1 != nil {
		log.Fatal(e1)
	}
}


func main() {
	var (
		addresserver string
		adressclient string
		e2    error
		server *snet.Addr

		udpConnection *snet.Conn
	)

	// Fetch arguments from command line
	flag.StringVar(&addresserver, "s", "", "Server SCION Address")
	flag.Parse()

	// Create the SCION UDP socket
	if len(addresserver) > 0 {
		server, err = snet.AddrFromString(addresserver)
		check(err)
	}

	dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(server.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, err = snet.ListenSCION("udp4", server)
	check(e2)

	receivePacketBuffer := make([]byte, 1000)//packet size
	for {
		n, adressclient, e2 := udpConnection.ReadFrom(receivePacketBuffer)
		time_recvd := time.Now().UnixNano()
		check(e2)

		_, size := binary.Uvarint(receivePacketBuffer)
		n = binary.PutVarint(receivePacketBuffer[size:], time_recvd)
		// Packet received, send back response to same client with time
		_, err = udpConnection.WriteTo(receivePacketBuffer[:n+size], addressclient)
		check(e2)
	}
}

