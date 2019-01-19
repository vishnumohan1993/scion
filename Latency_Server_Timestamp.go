
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

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func printUsage() {
	fmt.Println("\ntimestamp_server -s ServerSCIONAddress")
	fmt.Println("\tThe SCION address is specified as ISD-AS,[IP Address]:Port")
	//fmt.Println("\tListens for incoming connections and responds back to them right away")
	//fmt.Println("\tExample SCION address-s 19-ffaa:1:152,[192.168.0.102]:30102\n")
}

func main() {
	var (
		addressserver string

		err    error
		server *snet.Addr

		udpConnection *snet.Conn
	)


  flag.StringVar(&addressserver, "s", "", "Server SCION Address")
	flag.Parse()

  for len(addressserver) > 0 {
		server, err = snet.AddrFromString(addressserver)
		check(err)
	}// else {
		//printUsage()
		//check(fmt.Errorf("Error, server address needs to be specified with -s"))
	//}

	//for  len(addressserver)>0           //  statement refers to equatting  length of clientaddress with a condition
	//{
		//server, err = snet.AddrFromString(addressserver)
		//check(err)
	//}

	dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(server.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, err = snet.ListenSCION("udp4", server)
	check(err)

	bufferreceivePacket := make([]byte, 5000)//packet buffer array of size 5000 made for receiving
	for {
		n, addressclient, err := udpConnection.ReadFrom(bufferreceivePacket)
		check(err)

		// Packet received, send back response to same client
		m := binary.PutVarint(bufferreceivePacket[n:], time.Now().UnixNano())//packet received from source and encodes a uint64 into receivePacket and returns the number of bytes written and the func panic if buffer is too small
		_, err = udpConnection.WriteTo(bufferreceivePacket[: n+m], addressclient)
		check(err)
		fmt.Println("Received connection from", addressclient)
	}
}
