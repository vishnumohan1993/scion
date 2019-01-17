package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sciond"
)

func check(e error) { // Checking error loged in
	if e != nil {
		log.Fatal(e)
	}
}
func printUsage() {
	fmt.Println("\ndataplane_server -s ServerSCIONAddress")
	fmt.Println("\tListens for incoming connections and responds back to them right away")
	fmt.Println("\tThe SCION address is specified as ISD-AS,[IP Address]:Port")
	fmt.Println("\tExample SCION address 19-ffaa:1:14c,[192.168.0.110]:30100\n")
}
func main() {
	var (
		serverAddress string // variable for storing server address

		err    error
		server *snet.Addr

		udpConnection *snet.Conn
	)

	// Fetch arguments from command line
	flag.StringVar(&serverAddress, "s", "", "Server SCION Address")
	flag.Parse()

  if len(serverAddress) > 0 {   //// Create the SCION UDP socket
  		server, err = snet.AddrFromString(serverAddress) // adding server address for establishing connection
  		check(err)
  	} else {
  		printUsage()
  		check(fmt.Errorf("Error, server address needs to be specified with -s"))
  	}


    dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(server.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, err = snet.ListenSCION("udp4", server)
	check(err)

	receivePacketBuffer := make([]byte, 2500)
	for {
		n, clientAddress, err := udpConnection.ReadFrom(receivePacketBuffer)
		check(err)

		// Packet received, send back response to same client
		_, err = udpConnection.WriteTo(receivePacketBuffer[:n], clientAddress)
		check(err)
		fmt.Println("Received connection from", clientAddress)
	}
}
