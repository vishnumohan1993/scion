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

const (
	RECEIVE_SIZE int = 50000
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func printUsage() {
	fmt.Println("\nbw_est_server -s ServerSCIONAddress")
	fmt.Println("\tListens for incoming connections and responds back to them right away with the time received")
	fmt.Println("\tExample SCION address 19-ffaa:1:14c,[192.168.0.110]:30100)// our server AS
}

func main() {
	var (
		serverAddress string// variables

		err    error
		server *snet.Addr

		udpConnection *snet.Conn // Establishes Ip connection between server and clientAddress
	)

	// Fetch arguments from command line
	flag.StringVar(&serverAddress, "s", "", "Server SCION Address")
	flag.Parse()

	// Create the SCION UDP socket
	if len(serverAddress) > 0 {
		server, err = snet.AddrFromString(serverAddress) // adding server address to snet inorder to establish connections
		check(err)
	} else {
		printUsage()
		check(fmt.Errorf("Error, server address needs to be specified with -s"))
	}

	dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(server.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, err = snet.ListenSCION("udp4", server) // listen for scion client AS
	check(err)

	receivePacketBuffer := make([]byte, RECEIVE_SIZE + 1)
	for {
		n, clientAddress, err := udpConnection.ReadFrom(receivePacketBuffer)
		time_recvd := time.Now().UnixNano() // Taking time of packet receiving
		check(err)

		_, size := binary.Uvarint(receivePacketBuffer)
		n = binary.PutVarint(receivePacketBuffer[size:], time_recvd)
		// Packet received, send back response to same client with time
		_, err = udpConnection.WriteTo(receivePacketBuffer[:n+size], clientAddress)
		check(err)
	}
}
