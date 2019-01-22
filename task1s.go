
//reference example used from https://github.com/netsec-ethz/scion-homeworks/tree/master/latency
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


  flag.StringVar(&addressserver, "s", "", "Server SCION Address")
	flag.Parse()

	if len(addressserver)>0 {          //  statement refers to equatting  length of clientaddress with a condition
		server, e2 = snet.AddrFromString(addressserver)
		exceptioncheck(e2)
	}

	dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(server.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, e2 = snet.ListenSCION("udp4", server)
	exceptioncheck(e2)

	bufferreceivePacket := make([]byte, 5000)//packet buffer array of size 5000 made for receiving
	for {
		n, addressclient, e2 := udpConnection.ReadFrom(bufferreceivePacket)
		exceptioncheck(e2)

		// Packet received, send back response to same client
		m := binary.PutVarint(bufferreceivePacket[n:], time.Now().UnixNano())//packet received from source and encodes a uint64 into receivePacket and returns the number of bytes written and the func panic if buffer is too small
		_, e2 = udpConnection.WriteTo(bufferreceivePacket[: n+m], addressclient)
		exceptioncheck(e2)
		fmt.Println("Received connection from", addressclient)
	}
}
