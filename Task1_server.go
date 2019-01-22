//The code was developed by refering the hello world, sensorfetch app and master code of latency
// https://github.com/perrig/scionlab/blob/master/sensorapp/sensorfetcher/sensorfetcher.go
// https://github.com/netsec-ethz/scion-homeworks0/blob/master/latency/timestamp_server.go
// https://github.com/netsec-ethz/scion-apps/tree/master/helloworld

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func printUsage() {
	fmt.Println("\ntimestamp_server -s ServerSCIONAddress")
	fmt.Println("\tListens for incoming connections and responds back to them right away")
	fmt.Println("\tExample SCION address 1-1,[127.0.0.1]:42002\n")
}

func main() {
	var serverAddress string
	var err error
	var server *snet.Addr
	var udpConnection *snet.Conn

	// Fetch arguments from command line
	flag.StringVar(&serverAddress, "s", "", "Server SCION Address")
	flag.Parse()

	if len(serverAddress) == 0 {
		check(fmt.Errorf("Error, server address needs to be specified with -s"))
	}
	else {
	serverAddress, err = snet.AddrFromString(serverAddress)
	check(err)
	}
	
	dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(server.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, err = snet.ListenSCION("udp4", server)
	check(err)

	receivePacketBuffer := make([]byte, 2500)
	for {
		n, clientAddress, err := udpConnection.ReadFrom(receivePacketBuffer)
		check(err)
		time_recieved := time.Now().Unix()
		m := binary.PutVarint(receivePacketBuffer[n:], time_recieved)            //encoding the time of reciept to packet
		_, err = udpConnection.WriteTo(receivePacketBuffer[:n+m], clientAddress) //appending the recieved packet with the time stamp
		check(err)

		fmt.Println("Received connection from", clientAddress)
	}
}
