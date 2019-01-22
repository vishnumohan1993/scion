//The code was developed by refering the hello world, sensorfetch app and master code of latency 
// https://github.com/perrig/scionlab/blob/master/sensorapp/sensorfetcher/sensorfetcher.go
// https://github.com/netsec-ethz/scion-homeworks0/blob/master/latency/timestamp_server.go  
// https://github.com/netsec-ethz/scion-apps/tree/master/helloworld
package main

import (
	
	"flag"
	"fmt"
	"os"
	"log"
	"time"
	"math/rand"

	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/spath"
)

// Check just ensures the error is nil, or complains and quits
func Check(e error) {
	if e != nil {
		fmt.Println("Fatal error. Exiting.", "err", e)
		os.Exit(1)
	}
}
const (
	max_i = 10
	max_trails = 20
)
func printUsage() {
	fmt.Println(" -s ServerSCIONAddress -c ClientSCIONAddress")
	fmt.Println("The SCION address is specified as ISD-AS,[IP Address]:Port")
	fmt.Println("Example SCION address 1-1,[127.0.0.1]:42002")

func main() {
	var clientAddress string
	var serverAddress string
	var err error
	var local *snet.Addr
	var remote *snet.Addr
	var udpConnection *snet.Conn


	//Fetches the argument from command line
	flag.StringVar(&clientAddress, "c", "", "Client SCION Address")
	flag.StringVar(&serverAddress, "s", "", "Server SCION Address")
	flag.Parse()

	//SCION UDP socket creation
	
	if len(clientAddress) == 0
	{
		Check(fmt.Errorf("Error, local address needs to be specified with -local"))
	}
	else
	{
	clientAddress, err := snet.AddrFromString(clientAddress)
		Check(err)
	}
	if len(serverAddress) == 0 
	{
		Check(fmt.Errorf("Error, remote address needs to be specified with -remote"))
	}
	else
	{
		serverAddress, err := snet.AddrFromString(serverAddress)
		Check(err)
	}
	
	dispatcherAddr := "/run/shm/dispatcher/default.sock"

	snet.Init(local.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, err = snet.DialSCION("udp4", local, remote)
	check(err)


	random_seed := rand.NewSource(time.Now().Unix())
	
	var total int64 = 0
	i := 0 
	trials := 0 

	for i < max_i && trails < max_trails {
		num_tries += 1

		gen_id := rand.New(random_seed).Uint64()   //generating a random ID using rand function 
		n := binary.PutUvarint(sendPacketBuffer, gen_id)  //Using binary encoder and add the random number
		sendPacketBuffer[n] = 0

		time_sent := time.Now()
		_, err = udpConnection.Write(sendPacketBuffer)
		check(err)

		time_recieved := time.Now()
		_, _, err = udpConnection.ReadFrom(receivePacketBuffer)
		check(err)

		return_id, n := binary.Uvarint(receivePacketBuffer)
		
		if return_id == gen_id {
			time_received, _ := binary.Varint(receivePacketBuffer[n:]) // Using the binary encoder library, decoded the value is stored in time recieved 
			diff := (time_received.Unix() - time_sent.Unix())
			total += diff
			i += 1
		}
	}

	if iters != NUM_ITERS {
		check(fmt.Errorf("Error, exceeded maximum number of attempts"))
	}

	var final_diff float64 = float64(total) / float64(i)

	fmt.Printf("\nClient Address: %s\nServer Address: %s\n", clientAddress, serverAddress);
	fmt.Printf("RTT: %.3fs\n", final_diff)
	fmt.Printf( "Latency: %.3fs\n", final_diff)
}
