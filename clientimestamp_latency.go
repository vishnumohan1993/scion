package main

import (
	"flag"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"  //for mathematical computations
	"time"  // estimating latency time stamsps

	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sciond"
)


const (
	NUM_ITERS = 20
	MAX_NUM_TRIES = 40
)

func check(f error) {
	if f != nil {
		log.Fatal(f)
	}
}


func printUsage() {
	fmt.Println("\ntimestamp_client -c ClientSCIONAddress -s ServerSCIONAddress")
	fmt.Println("\tProvides speed estimates (RTT and latency) from client to dedicated server")
	fmt.Println("\tThe SCION address is specified as ISD-AS,[IP Address]:Port")
	fmt.Println("\tExample SCION address  -s 19-ffaa:1:152,[192.168.0.102]:30102\n") //an AS Id
}


func main() {

	var (
		clientAddress string  //local variables
	  serverAddress string

		err    error
		local  *snet.Addr
		remote *snet.Addr

		udpConnection *snet.Conn
	)

	flag.StringVar(&clientAddress, "c", "", "Client SCION Address")
	flag.StringVar(&serverAddress, "s", "", "Server SCION Address")
	flag.Parse()

	if len(clientAddress) > 0 {           //creating Scion UDP socket
		local, err = snet.AddrFromString(clientAddress)
		check(err) //passing the error to check function for loging it

		} else {
		printUsage()
		check(fmt.Errorf("Error, client address needs to start with -c"))
	}

	if len(serverAddress) > 0 {
		remote, err = snet.AddrFromString(serverAddress)   //Adding server address for establishing connection
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

	var total int64 = 0
	iters := 0
	num_tries := 0
	for iters < NUM_ITERS && num_tries < MAX_NUM_TRIES {
		num_tries += 1

		id := rand.New(seed).Uint64()
		n := binary.PutUvarint(sendPacketBuffer, id)
		sendPacketBuffer[n] = 0

		time_sent := time.Now()
		_, err = udpConnection.Write(sendPacketBuffer)
		check(err)

		_, _, err = udpConnection.ReadFrom(receivePacketBuffer)
		check(err)

		ret_id, n := binary.Uvarint(receivePacketBuffer)
		if ret_id == id {
			time_received, _ := binary.Varint(receivePacketBuffer[n:])
			diff := (time_received - time_sent.UnixNano())
			total += diff
			iters += 1
		}
	}

	if iters != NUM_ITERS {
		check(fmt.Errorf("Error, exceeded maximum number of attempts"))
	}

	var difference float64 = float64(total) / float64(iters)

	fmt.Printf("\nClient: %s\nServer: %s\n", clientAddress, serverAddress);
	fmt.Println("Time estimates:")
	// Print in ms, so divide by 1e6 from nano
	fmt.Printf("\tRTT - %.3fms\n", difference/1e6)
	fmt.Printf("\tLatency - %.3fms\n", difference/2e6)
}
