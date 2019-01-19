// Example reference used from https://github.com/netsec-ethz/scion-homeworks/tree/master/latency

package main

import (
	"flag"
	"encoding/binary"
	"fmt"
	//"log"
	"math/rand"  //for mathematical computations
	"time"  // estimating latency time stamps

	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sciond"
)




//func check(f error) {
	//if f != nil {
	//	log.Fatal(f)
	//}
//}

const (
	count = 15 // iteration count
	limit = 30// maximum allowed of above count
)


func printUsage() {
	fmt.Println("\ntimestamp_client -c ClientSCIONAddress -s ServerSCIONAddress")
	fmt.Println("\tThe SCION address is specified as ISD-AS,[IP Address]:Port")
}


func main() {

	var (
		addressclient string  
	    addressserver string

		err    error
		local  *snet.Addr
		remote *snet.Addr

		udpConnection *snet.Conn
	)

	flag.StringVar(&addressclient, "c", "", "Client SCION Address")  // binding the flag to client address string
	flag.StringVar(&addressserver, "s", "", "Server SCION Address") // binding the flag to server address string
	flag.Parse() // parsing command line to above defined flags


//Creating Scion UDP socket
	//if  len(addressclient)>0           //  statement refers to equatting  length of clientaddress with a condition
	//{
		//local, err = snet.AddrFromString(addressclient)
		//check(err)
	//}

	for len(addressclient) > 0 {           
		local, err = snet.AddrFromString(addressclient)
		//check(err) //passing the error to check function for logging it

		} 
		//else
		//{
		//printUsage()
		//check(fmt.Errorf("Error, client address needs to start with -c"))
	//}

	//if  len(addressserver)>0           //  statement refers to equatting  length of clientaddress with a condition
	//{
		//local, err = snet.AddrFromString(addressserver)
		//check(err)
	//}

	for len(addressserver) > 0 { // same as above
		remote, err = snet.AddrFromString(addressserver)   //Adding server address for establishing connection
		//check(err)
	} 
	//else {
		//printUsage()
		//check(fmt.Errorf("Error, server address needs to be specified with -s"))
	//}

dispatcherAddr := "/run/shm/dispatcher/default.sock"
snet.Init(local.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

udpConnection, err = snet.DialSCION("udp4", local, remote)
	//check(err)

	bufferreceivePacket := make([]byte, 5000) //packet buffer array of size 5000 made for receiving
	buffersendPacket := make([]byte, 25) //packet buffer array of size 30 made for sending

	seed := rand.NewSource(time.Now().UnixNano())  //creating seed with random ids from new source


	//calculation starts here

	var final int64 = 0

	i := 0 //iterations
	j := 0 //limit

	for i < count && j < limit {
		j += 1

		id := rand.New(seed).Uint64() // random number generation with validation when destination sends packet back

		n := binary.PutUvarint(buffersendPacket, id) //encodes a uint64 into sendPacket and returns the number of bytes written and the func panic if buffer is too small
		buffersendPacket[n] = 0

		t1 := time.Now() //fetching the time sent for later calculations
		_, err = udpConnection.Write(buffersendPacket)
		//check(err)

		_, _, err = udpConnection.ReadFrom(bufferreceivePacket)
		//check(err)

		ret_id, n := binary.Uvarint(bufferreceivePacket)
		//decoding
		if ret_id == id {
			t2, _ := binary.Varint(bufferreceivePacket[n:])//time of receive
			difference := (t2 - t1.UnixNano())  //unixnano refers to nanosecond range of time
			total+ = difference
			i+=1
		}
	}

	//if i != count {
	//	check(fmt.Errorf("Error, exceeded maximum number of attempts"))
	//}

	var difference2 float64 = float64(total) / float64(i)//final result

	fmt.Printf("\nClient Address, IP and Port: %s\n", addressclient);
	fmt.Printf("\nServer Address, IP and Port: %s\n",addressserver);
	fmt.Println("Results Obtained as follows:")
	// Result is printed in milliseconds, so divide by 1e6 from nano
	fmt.Printf("\tRTT - %.3fms\n", difference2/1e6)
	fmt.Printf("\tLatency - %.3fms\n", difference2/2e6)//since we take RTT as 2x Latency
}
