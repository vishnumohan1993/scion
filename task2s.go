//reference example used from https://github.com/netsec-ethz/scion-homeworks/tree/master/latency

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sciond"
)

// to display incase of any issues during runtime
func exceptioncheck(e1 error) {
	if e1 != nil {
		log.Fatal(e1)
	}
}

func main() {
	var (
		addressserver string
		e2    error
		server *snet.Addr
		udpConnection *snet.Conn
	)

	// Fetch arguments from command line
	flag.StringVar(&addressserver, "s", "", "Server SCION Address")
	flag.Parse()

	if len(addressserver)>0 {           //  statement refers to equatting  length of clientaddress with a condition
		server, e2 = snet.AddrFromString(addressserver)
		exceptioncheck(e2)
	}

    dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(server.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, e2 = snet.ListenSCION("udp4", server)
	exceptioncheck(e2)

	bufferreceivePacket := make([]byte, 5000)
	for {
		n, addressclient, e2 := udpConnection.ReadFrom(bufferreceivePacket)
		exceptioncheck(e2)

		// Packet received, send back response to same client
		_, e2 = udpConnection.WriteTo(bufferreceivePacket[:n], addressclient)
		exceptioncheck(e2)
		fmt.Println("Received connection from", addressclient)
	}
}
