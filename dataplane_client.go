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
const (
	NUM_ITERS = 20
	MAX_NUM_TRIES = 40
)

func check(e error) { // Checking error loged in
	if e != nil {
		log.Fatal(e)
	}
}
func printUsage() {
	fmt.Println("\ndataplane_client -c clientSCIONAddress -s serverSCIONAddress")
	fmt.Println("\tProvides speed estimates (RTT and latency) from client to dedicated server")
	fmt.Println("\tExample SCION address 19-ffaa:1:14c,[192.168.0.110]:30100\n")
}
func main() {
	var (
		clientAddress string //variables for adress values
		serverAddress string

		err    error
		local  *snet.Addr
		remote *snet.Addr

		udpConnection *snet.Conn
	)

  flag.StringVar(&clientAddress, "c", "", "client SCION Address")
  	flag.StringVar(&serverAddress, "s", "", "server SCION Address")
  	flag.Parse()

    if len(clientAddress) > 0 { Create the SCION UDP socket
		local, err = snet.AddrFromString(serverAddress) // adding server address for connection establishment
		check(err)
	} else {
		printUsage()
		check(fmt.Errorf("Error, client address needs to be specified with -c"))
	}
  if len(serverAddress) > 0 {
		remote, err = snet.AddrFromString(serverAddress)
		check(err)
	} else {
		printUsage()
		check(fmt.Errorf("Error, server address needs to be specified with -s"))
	}

  dispatcherAddr := "/run/shm/dispatcher/default.sock"
snet.Init(local.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

udpConnection, err = snet.DialSCION("udp4", local, remote)
check(err)

receivePacketBuffer := make([]byte, 2500)
sendPacketBuffer := make([]byte, 16)

seed := rand.NewSource(time.Now().UnixNano())
// Do 5 iterations so we can use average
var total int64 = 0
iters := 0
num_tries := 0
for iters < NUM_ITERS && num_tries < MAX_NUM_TRIES {
  num_tries += 1

  id := rand.New(seed).Uint64()
  n := binary.PutUvarint(sendPacketBuffer, id)
  sendPacketBuffer[n] = 0

  time_sent := time.Now() // taking timestamp now
  _, err = udpConnection.Write(sendPacketBuffer)
  check(err)

  _, _, err = udpConnection.ReadFrom(receivePacketBuffer) //read the received packet
  time_received := time.Now()
  check(err)

  ret_id, n := binary.Uvarint(receivePacketBuffer)
  if ret_id == id {
    diff := (time_received.UnixNano() - time_sent.UnixNano()) // taking the difference value 
    total += diff
    iters += 1
    // fmt.Printf("%d: %.3fms %.3fms\n", iters, float64(diff)/1e6, float64(diff)/2e6)
  }
}

if iters != NUM_ITERS {
  check(fmt.Errorf("Error, exceeded maximum number of attempts"))
}

var difference float64 = float64(total) / float64(iters)

fmt.Printf("\nSource: %s\nDestination: %s\n", clientAddress, serverAddress);
fmt.Println("Time estimates:")
// Print in ms, so divide by 1e6 from nano
fmt.Printf("\tRTT - %.3fms\n", difference/1e6)
fmt.Printf("\tLatency - %.3fms\n", difference/2e6)
}
