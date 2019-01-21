//reference example used from https://github.com/netsec-ethz/scion-homeworks/tree/master/latency

package main

import (
	"flag"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
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
		addressclient string  
	    	addressserver string

		e2    error
		local  *snet.Addr
		remote *snet.Addr

		udpConnection *snet.Conn
	)

  	flag.StringVar(&addressclient, "c", "", "client SCION Address")
  	flag.StringVar(&addressserver, "s", "", "server SCION Address")
  	flag.Parse()


  	if len(addressclient)>0 {          //  statement refers to equatting  length of clientaddress with a condition
		local, e2 = snet.AddrFromString(addressclient)
		 exceptioncheck(e2)
	}


	if len(addressserver)>0 {           //  statement refers to equatting  length of clientaddress with a condition
		local, e2 = snet.AddrFromString(addressserver)
		exceptioncheck(e2)
	}

 dispatcherAddr := "/run/shm/dispatcher/default.sock"
snet.Init(local.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

udpConnection, e2 = snet.DialSCION("udp4", local, remote)
exceptioncheck(e2)

bufferreceivePacket := make([]byte, 5000)//packet buffer array of size 5000 made for receiving
buffersendPacket := make([]byte, 25)//packet buffer array of size 30 made for sending

seed := rand.NewSource(time.Now().UnixNano())//creating seed with random ids from new source


  id := rand.New(seed).Uint64()// random number generation with validation when destination sends packet back

  n := binary.PutUvarint(buffersendPacket, id)//encodes a uint64 into sendPacket and returns the number of bytes written and the func panic if buffer is too small
  buffersendPacket[n] = 0

  t1 := time.Now() //fetching the time at which packet is sent
  _, e2 = udpConnection.Write(buffersendPacket)
  exceptioncheck(e2)

  _, _, e2 = udpConnection.ReadFrom(bufferreceivePacket) 
  t2 := time.Now()//fetching time when packet is back to origin
  exceptioncheck(e2)

  ret_id, n := binary.Uvarint(bufferreceivePacket)
  if ret_id == id {
   var difference float64 = (float64(t2.UnixNano()) - float64(t1.UnixNano())) // unixnano refers to nanosecond range of time 
  }
fmt.Printf("\nClient Address, IP and Port: %s\n",addressclient)
fmt.Printf("\nServer Address, IP and Port: %s\n",addressserver)
fmt.Println("Results Obtained as follows:")

// Result is printed in milliseconds, so divide by 1e6 from nano
var difference float64
	
fmt.Printf("\tRTT - %.3fms\n", difference/1e6)
fmt.Printf("\tLatency - %.3fms\n", difference/2e6)//since we take RTT as 2x Latency
}
