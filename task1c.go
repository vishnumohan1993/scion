                                                                                                                     // Example reference used from https://github.com/netsec-ethz/scion-homeworks/tree/master/latency
//reference example used from https://github.com/netsec-ethz/scion-homeworks/tree/master/latency
package main

import (
	"flag"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"  //for mathematical computations
	"time"  // estimating latency time stamps

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
//required variable declarations with data types
	var (
		addressclient string  
	    	addressserver string

		e2    error
		local  *snet.Addr
		remote *snet.Addr

		udpConnection *snet.Conn
	)

	flag.StringVar(&addressclient, "c", "", "Client SCION Address")  // binding the flag to client address string
	flag.StringVar(&addressserver, "s", "", "Server SCION Address") // binding the flag to server address string
	flag.Parse() // parsing command line to above defined flags


//Creating Scion UDP socket
	if  len(addressclient) > 0 {          //  statement refers to equatting  length of clientaddress with a condition
		local, e2 = snet.AddrFromString(addressclient)
		 exceptioncheck(e2)
	}

	if  len(addressserver) > 0 {           //  statement refers to equatting  length of serveraddress with a condition
		remote, e2 = snet.AddrFromString(addressserver)
		exceptioncheck(e2)
	}


dispatcherAddr := "/run/shm/dispatcher/default.sock"
snet.Init(local.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

udpConnection,e2 = snet.DialSCION("udp4", local, remote)
	exceptioncheck(e2)

	bufferreceivePacket := make([]byte, 5000) //packet buffer array of size 5000 made for receiving
	buffersendPacket := make([]byte, 25) //packet buffer array of size 30 made for sending

	seed := rand.NewSource(time.Now().UnixNano())  //creating seed with random ids from new source
		var ans int64 = 0

		id := rand.New(seed).Uint64() // random number generation with validation when destination sends packet back

		n := binary.PutUvarint(buffersendPacket, id) //encodes a uint64 into sendPacket and returns the number of bytes written and the func panic if buffer is too small
		buffersendPacket[n] = 0//passing the value to the same

		t1 := time.Now() //fetching the time at which packet is sent
		_, e2 = udpConnection.Write(buffersendPacket)
		exceptioncheck(e2)

		_, _, e2 = udpConnection.ReadFrom(bufferreceivePacket)
		exceptioncheck(e2)

		ret_id, n := binary.Uvarint(bufferreceivePacket)//calculation starts when the id which retuens back matches to initail sent id
		// final result is plotted when returning id matches with original one
		if ret_id == id {
			t2, _ := binary.Varint(bufferreceivePacket[n:])//time of receive
			difference := (t2 - t1.UnixNano())  //unixnano refers to nanosecond range of time 
			ans = difference
		}
	var result float64 = float64(ans)
	//fmt.Printf("\nClient Address, IP and Port: %s\n",addressclient);
	fmt.Printf("\nSource: %s\nDestination: %s\n", addressclient, addressserver);
	//fmt.Printf("\nServer Address, IP and Port: %s\n",addressserver);
	fmt.Println("Results Obtained as follows:")

	// Result is printed in milliseconds, so divide by 1e6 from nano

	fmt.Printf("\tRTT - %.3fms\n", result/1e6)
	fmt.Printf("\tLatency - %.3fms\n", result/2e6)//since we take RTT as 2x Latency
}
